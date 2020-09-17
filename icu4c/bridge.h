#ifndef __C_BRIDGE_H__
#define __C_BRIDGE_H__

#include <stdlib.h>
#include <stdbool.h>
#include <unicode/utypes.h>
#include <unicode/udat.h>

const UErrorCode go_format_datetime(
	const char* locale,
	const char* tz,
	UDateFormatStyle date_style,
	UDateFormatStyle time_style,
	double msec,
	char* const result,
	const size_t result_size
);

#endif //__C_BRIDGE_H__
