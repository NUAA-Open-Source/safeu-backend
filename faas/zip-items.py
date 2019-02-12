# Code for Python 2.7

import os
import oss2
import zipfile
import shutil
import time
import json
import urllib

oss_endpoint = "oss-cn-shanghai.aliyuncs.com"
oss_bucket_name = "safeu"

def handler(environ, start_response):
    context = environ['fc.context']
    request_uri = environ['fc.request_uri']
    for k, v in environ.items():
        if k.startswith("HTTP_"):
            # process custom request headers
            print(k, v)
            pass

    # get request body
    try:
        request_body_size = int(environ.get('CONTENT_LENGTH', 0))
    except(ValueError):
        request_body_size = 0
    request_body = environ['wsgi.input'].read(request_body_size)
    request_body = urllib.unquote(request_body).decode('utf8') 

    # get request method
    request_method = environ['REQUEST_METHOD']
    if request_method != 'POST':
        status = '400 Bad Request'
        response_headers = [('Content-type', 'application/json')]
        start_response(status, response_headers)
        data = json.dumps({"error": "invalid request method."})
        return [data]

    # print request body
    print('request_body: {}'.format(request_body))
    request_body_json = json.loads(request_body)

    creds = context.credentials
    auth = oss2.StsAuth(creds.accessKeyId, creds.accessKeySecret, creds.securityToken)
    bucket = oss2.Bucket(auth, oss_endpoint, oss_bucket_name)
    #your source list
    # sourceFile = ['resource/1.jpg','resource/2.jpg']
    sourceFile = request_body_json.get("items")

    #zip name
    uid = request_body_json.get("re_code")
    tmpdir = '/tmp/download/'

    os.system("rm -rf /tmp/*")
    os.mkdir(tmpdir)

    #download
    for name in sourceFile :
        millis = int(round(time.time() * 1000))
        bucket.get_object_to_file(name , tmpdir + name)

    #zip file
    zipname = '/tmp/'+uid + '.zip'
    make_zip(tmpdir , zipname)

    #upload
    total_size = os.path.getsize(zipname)
    part_size = oss2.determine_part_size(total_size, preferred_size = 128 * 1024)

    key = 'archive/' + uid + '.zip'
    upload_id = bucket.init_multipart_upload(key).upload_id

    with open(zipname, 'rb') as fileobj:
        parts = []
        part_number = 1
        offset = 0
        while offset < total_size:
            num_to_upload = min(part_size, total_size - offset)
            result = bucket.upload_part(key, upload_id, part_number,oss2.SizedFileAdapter(fileobj, num_to_upload))
            parts.append(oss2.models.PartInfo(part_number, result.etag))
            offset += num_to_upload
            part_number += 1

        bucket.complete_multipart_upload(key, upload_id, parts)
    
    status = '200 OK'
    response_headers = [('Content-type', 'application/json')]
    start_response(status, response_headers)
    url = "https://" + oss_bucket_name + "." + oss_endpoint + "/archive/" + uid + ".zip"
    data = json.dumps({"url": url})
    return [data]


def make_zip(source_dir, output_filename):
    zipf = zipfile.ZipFile(output_filename, 'w')
    pre_len = len(os.path.dirname(source_dir))
    for parent, dirnames, filenames in os.walk(source_dir):
        for filename in filenames:
            pathfile = os.path.join(parent, filename)
            arcname = pathfile[pre_len:].strip(os.path.sep)
            zipf.write(pathfile, arcname)
    zipf.close()
