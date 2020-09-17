#include <string.h>
#include <unicode/ustring.h>
#include <unicode/udat.h>

#include "bridge.h"

const UErrorCode go_format_datetime(
	const char* locale,
	const char* tz,
	UDateFormatStyle date_style,
	UDateFormatStyle time_style,
	double msec,
	char* const result,
	const size_t result_size
) {
	UErrorCode status = U_ZERO_ERROR;
	UChar buf[result_size];

	UChar tzUchar[strlen(tz) * sizeof(UChar)];
	u_uastrcpy(tzUchar, tz);

	UDateFormat* fmt = udat_open(
		time_style,
		date_style,
		locale,
		tzUchar,
		-1, // -1 because tzUchar is null-terminated.
		NULL, // pattern is unused here.
		-1, // pattern is unused here.
		&status
	);
	if (U_FAILURE(status)) {
		goto exit0;
	}

	udat_format(
		fmt,
		(UDate)msec,
		buf,
		result_size,
		NULL,
		&status
	);
	if (U_FAILURE(status)) {
		goto exit1;
	}

	u_austrcpy(result, buf);
exit1:
	udat_close(fmt);
exit0:
	return status;
}
