package common

var Errors = map[int]string{
	0:  "General error.",
	1:  "Parameter error.",
	2:  "The length of the time is not within the right range.",
	3:  "Can't Find User Token.",
	4:  "Repeat retrieve code.",
	5:  "Auth field required.",
	6:  "The retrieve code mismatch auth.",
	7:  "Cannot find resource via this retrieve code.",
	8:  "There is a problem for this resource, please contact the maintainer.",
	9:  "Cannot get the password.",
	10: "The password is not correct.",
	11: "Cannot get the token.",
	12: "Token invalid.",
	13: "Over the expired time.",
	14: "Out of downloadable count.",
	15: "Cannot get the items field.",
	16: "Cannot get the download link.",
	17: "Service Unavailable, please contact the maintainer.",
	18: "Illegal callback request in Upload",
}
