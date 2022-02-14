package utils

// NormalizePhoneNumber normalized phone number based on E.164 standard
func NormalizePhoneNumber(phoneNumber string) string {
	if len(phoneNumber) < 4 {
		return ""
	}

	var areaPhoneNumbersIndonesian = "+62"
	var convertedNumber string
	if phoneNumber[:2] == "08" {
		convertedNumber = areaPhoneNumbersIndonesian + phoneNumber[1:]
	} else if phoneNumber[:1] == "8" {
		convertedNumber = areaPhoneNumbersIndonesian + phoneNumber
	} else if phoneNumber[:3] == "628" {
		convertedNumber = areaPhoneNumbersIndonesian + phoneNumber[2:]
	} else if phoneNumber[:3] == areaPhoneNumbersIndonesian {
		convertedNumber = phoneNumber
	}
	return convertedNumber
}

// NormalizePhoneNumberLocal normalized phone number based on local standard
func NormalizePhoneNumberLocal(phoneNumber string) string {
	if len(phoneNumber) < 4 {
		return ""
	}

	var areaPhoneNumbersIndonesian = "0"
	var convertedNumber string
	if phoneNumber[:4] == "+628" {
		convertedNumber = areaPhoneNumbersIndonesian + phoneNumber[3:]
	} else if phoneNumber[:1] == "8" {
		convertedNumber = areaPhoneNumbersIndonesian + phoneNumber
	} else if phoneNumber[:3] == "628" {
		convertedNumber = areaPhoneNumbersIndonesian + phoneNumber[2:]
	} else if phoneNumber[:3] == areaPhoneNumbersIndonesian {
		convertedNumber = phoneNumber
	}
	return convertedNumber
}

func GetCodePhoneNumber(phoneNumber string) string {
	nor := NormalizePhoneNumber(phoneNumber)
	return "08" + nor[4:6]
}
