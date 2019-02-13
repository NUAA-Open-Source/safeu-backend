# Code for Python 2.7

import os
import oss2
import zipfile
import shutil
import time
import json
import urllib

# oss_endpoint = "oss-cn-shanghai.aliyuncs.com"
# oss_bucket_name = "safeu"

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

    items = request_body_json.get("items")
    print("[DEBUG] items: {0}".format(items))

    # zip name
    re_code = request_body_json.get("re_code")
    is_full = request_body_json.get("full")
    tmpdir = '/tmp/download/'

    os.system("rm -rf /tmp/*")
    os.mkdir(tmpdir)

    # download
    for item in items :
        print("[DEBUG] item: {}".format(item))

        oss_protocol = item.get("protocol")
        oss_bucket_name = item.get("bucket")
        oss_endpoint = item.get("endpoint")
        file_path = item.get("path")
        file_original_name = item.get("original_name")

        bucket = oss2.Bucket(auth, oss_endpoint, oss_bucket_name)
        
        bucket.get_object_to_file(file_path , tmpdir + file_original_name)

    #zip file
    zipname = '/tmp/'+ re_code + '.zip'
    make_zip(tmpdir , zipname)

    #upload
    total_size = os.path.getsize(zipname)
    part_size = oss2.determine_part_size(total_size, preferred_size = 128 * 1024)

    if is_full:
        zip_path = 'full-archive/' + re_code + '.zip'
    else:
        zip_path = 'custom-archive/' + re_code + '.zip'

    # use the last bucket to upload zip package
    upload_id = bucket.init_multipart_upload(zip_path).upload_id

    with open(zipname, 'rb') as fileobj:
        parts = []
        part_number = 1
        offset = 0
        while offset < total_size:
            num_to_upload = min(part_size, total_size - offset)
            result = bucket.upload_part(zip_path, upload_id, part_number,oss2.SizedFileAdapter(fileobj, num_to_upload))
            parts.append(oss2.models.PartInfo(part_number, result.etag))
            offset += num_to_upload
            part_number += 1

        bucket.complete_multipart_upload(zip_path, upload_id, parts)

    zip_meta = bucket.head_object(zip_path)
    zip_content_type = zip_meta.headers.get('Content-Type')
    
    status = '200 OK'
    response_headers = [('Content-type', 'application/json')]
    start_response(status, response_headers)
    url = "https://" + oss_bucket_name + "." + oss_endpoint + "/" + zip_path
    data = json.dumps({
        "host": url,
        "protocol": oss_protocol,
        "bucket": oss_bucket_name,
        "endpoint": oss_endpoint,
        "path": zip_path,
        "original_name": re_code + ".zip",
        "type": zip_content_type,
    })
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
