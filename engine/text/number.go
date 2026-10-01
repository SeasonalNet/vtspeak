package text

import (
	"errors"
	"fmt"
	"strings"
)

var cardinalOnes = [...]string{
	"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine",
	"ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen",
	"seventeen", "eighteen", "nineteen",
}

var cardinalTens = [...]string{
	"", "", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety",
}

var cardinalScales = [...]string{"", "thousand", "million", "billion", "trillion"}

var ordinalUnits = [...]string{
	"zeroth", "first", "second", "third", "fourth", "fifth", "sixth", "seventh", "eighth", "ninth",
	"tenth", "eleventh", "twelfth", "thirteenth", "fourteenth", "fifteenth", "sixteenth",
	"seventeenth", "eighteenth", "nineteenth",
}

var ordinalTens = [...]string{"", "", "twentieth", "thirtieth"}

// ExpandPaul2013UnsignedInteger spells one unpunctuated unsigned integer. It
// uses cardinal groups through 15 digits, applies the observed four-digit year
// rule, and spells longer values or multi-digit leading-zero values digit by
// digit.
func ExpandPaul2013UnsignedInteger(value string) ([]string, error) {
	if value == "" {
		return nil, errors.New("integer surface is empty")
	}
	for position := 0; position < len(value); position++ {
		if value[position] < '0' || value[position] > '9' {
			return nil, fmt.Errorf("integer surface contains non-digit byte at position %d", position)
		}
	}
	if value == "0" {
		return []string{"zero"}, nil
	}
	if len(value) > 15 || len(value) > 1 && value[0] == '0' {
		words := make([]string, len(value))
		for index := range value {
			if value[index] == '0' {
				words[index] = "oh"
			} else {
				words[index] = cardinalOnes[value[index]-'0']
			}
		}
		return words, nil
	}
	if len(value) == 4 && value[0] != '0' {
		return expandFourDigitYear(value), nil
	}
	return expandCardinalGroups(value), nil
}

func expandCardinalGroups(value string) []string {
	words := make([]string, 0, len(value))
	groups := 0
	for end := len(value); end > 0; end -= 3 {
		start := end - 3
		if start < 0 {
			start = 0
		}
		group := 0
		for index := start; index < end; index++ {
			group = group*10 + int(value[index]-'0')
		}
		if group != 0 {
			part := cardinalBelowThousand(group)
			if cardinalScales[groups] != "" {
				part = append(part, cardinalScales[groups])
			}
			words = append(part, words...)
		}
		groups++
	}
	return words
}

// ExpandPaul2013Number spells the captured plain integer, signed integer,
// and decimal forms. Grouping commas are accepted only in three-digit groups.
// Decimal fractional digits are spoken individually. Dates, times, currency,
// percentages, slash dates, and ordinal suffixes are handled by separate
// normalizers.
func ExpandPaul2013Number(value string) ([]string, error) {
	if value == "" {
		return nil, errors.New("number surface is empty")
	}
	var sign string
	if value[0] == '+' || value[0] == '-' {
		if value[0] == '+' {
			sign = "plus"
		} else {
			sign = "minus"
		}
		value = value[1:]
		if value == "" {
			return nil, errors.New("number sign has no digits")
		}
	}

	if strings.Count(value, ".") > 1 {
		return nil, errors.New("number contains more than one decimal point")
	}
	integerPart, fractionalPart, hasDecimal := strings.Cut(value, ".")
	if hasDecimal && fractionalPart == "" {
		return nil, errors.New("decimal point has no fractional digits")
	}
	if integerPart == "" && !hasDecimal {
		return nil, errors.New("number has no integer digits")
	}

	words := make([]string, 0)
	if integerPart != "" {
		ungrouped, err := ungroupInteger(integerPart)
		if err != nil {
			return nil, err
		}
		words, err = ExpandPaul2013UnsignedInteger(ungrouped)
		if err != nil {
			return nil, err
		}
	} else if sign != "" {
		return nil, errors.New("signed decimal requires an integer part")
	}
	if hasDecimal {
		words = append(words, "point")
		for position := 0; position < len(fractionalPart); position++ {
			digit := fractionalPart[position]
			if digit < '0' || digit > '9' {
				return nil, fmt.Errorf("fraction contains non-digit byte at position %d", position)
			}
			words = append(words, cardinalOnes[digit-'0'])
		}
	}
	if sign != "" {
		words = append([]string{sign}, words...)
	}
	return words, nil
}

// ExpandPaul2013Currency handles a leading dollar sign and an optional
// two-digit cents field. The captured `$5.00` form becomes “five dollars”.
// Singular/plural, nonzero-cent wording, and other amounts are implementation
// rules inferred from that captured form, not claims of DLL parity.
func ExpandPaul2013Currency(value string) ([]string, error) {
	if !strings.HasPrefix(value, "$") {
		return nil, errors.New("currency surface must begin with $")
	}
	amount := value[1:]
	if amount == "" || strings.HasPrefix(amount, "+") || strings.HasPrefix(amount, "-") {
		return nil, errors.New("currency amount must be unsigned and nonempty")
	}
	if strings.Count(amount, ".") > 1 {
		return nil, errors.New("currency amount contains more than one decimal point")
	}
	dollars, cents, hasDecimal := strings.Cut(amount, ".")
	if hasDecimal && len(cents) != 2 {
		return nil, errors.New("currency cents field must contain exactly two digits")
	}
	if dollars == "" {
		return nil, errors.New("currency amount requires a dollar field")
	}
	ungroupedDollars, err := ungroupInteger(dollars)
	if err != nil {
		return nil, fmt.Errorf("currency dollar field: %w", err)
	}
	if len(ungroupedDollars) > 15 {
		return nil, errors.New("currency dollar field exceeds the supported 15-digit range")
	}
	dollarWords, err := ExpandPaul2013UnsignedInteger(ungroupedDollars)
	if err != nil {
		return nil, err
	}
	words := append([]string(nil), dollarWords...)
	if ungroupedDollars == "1" {
		words = append(words, "dollar")
	} else {
		words = append(words, "dollars")
	}
	if hasDecimal {
		if err := validateDigits(cents, "currency cents"); err != nil {
			return nil, err
		}
		centValue := int(cents[0]-'0')*10 + int(cents[1]-'0')
		if centValue != 0 {
			words = append(words, cardinalBelowThousand(centValue)...)
			if centValue == 1 {
				words = append(words, "cent")
			} else {
				words = append(words, "cents")
			}
		}
	}
	return words, nil
}

// ExpandPaul2013Percentage appends the captured percentage marker wording to
// an ordinary supported integer or decimal number. Only `25%` is directly
// represented in the local runtime captures; the broader numeric grammar is
// shared with ExpandPaul2013Number and remains a parity hypothesis.
func ExpandPaul2013Percentage(value string) ([]string, error) {
	if !strings.HasSuffix(value, "%") {
		return nil, errors.New("percentage surface must end with %")
	}
	number := value[:len(value)-1]
	words, err := ExpandPaul2013Number(number)
	if err != nil {
		return nil, fmt.Errorf("percentage value: %w", err)
	}
	return append(words, "percent"), nil
}

// ExpandPaul2013Telephone spells the captured NNN-NNNN surface as cardinal
// groups separated by "to". The Stage 20 parser-row trace for `555-1234`
// directly shows "five hundred fifty five to twelve thirty four". Other
// digit values within this layout use the existing cardinal expander as an
// implementation inference; other telephone layouts are unsupported.
func ExpandPaul2013Telephone(value string) ([]string, error) {
	parts := strings.Split(value, "-")
	if len(parts) != 2 || len(parts[0]) != 3 || len(parts[1]) != 4 {
		return nil, errors.New("telephone surface must use the captured NNN-NNNN form")
	}
	for index, part := range parts {
		if err := validateDigits(part, fmt.Sprintf("telephone group %d", index+1)); err != nil {
			return nil, err
		}
	}
	left, err := ExpandPaul2013UnsignedInteger(parts[0])
	if err != nil {
		return nil, fmt.Errorf("telephone first group: %w", err)
	}
	right, err := ExpandPaul2013UnsignedInteger(parts[1])
	if err != nil {
		return nil, fmt.Errorf("telephone second group: %w", err)
	}
	words := append([]string(nil), left...)
	words = append(words, "to")
	return append(words, right...), nil
}

// ExpandPaul2013Ordinal expands ordinal digit surfaces through thirty-first.
// The observed normalizer uses a dedicated suffix-sensitive lookup path and
// has a special zero entry. The 1–31 spellings here follow English ordinal
// morphology; only 2nd (in a captured date) and the word fifth (from a
// captured date) have direct runtime output evidence, so the remaining forms
// are implementation coverage rather than claimed DLL parity.
func ExpandPaul2013Ordinal(value string) ([]string, error) {
	if len(value) < 3 {
		return nil, errors.New("ordinal surface must contain digits and a suffix")
	}
	suffix := value[len(value)-2:]
	if suffix != "st" && suffix != "nd" && suffix != "rd" && suffix != "th" {
		return nil, errors.New("ordinal surface has an unsupported suffix")
	}
	digits := value[:len(value)-2]
	if err := validateDigits(digits, "ordinal"); err != nil {
		return nil, err
	}
	if len(digits) > 2 {
		return nil, errors.New("ordinal value is outside the supported range 0 through 31")
	}
	if len(digits) > 1 && digits[0] == '0' {
		return nil, errors.New("ordinal surface with a leading zero is unsupported")
	}
	valueNumber := 0
	for position := 0; position < len(digits); position++ {
		valueNumber = valueNumber*10 + int(digits[position]-'0')
	}
	if valueNumber > 31 {
		return nil, errors.New("ordinal value is outside the supported range 0 through 31")
	}
	if expected := ordinalSuffix(valueNumber); suffix != expected {
		return nil, fmt.Errorf("ordinal suffix %q does not match value %d", suffix, valueNumber)
	}
	return expandOrdinalValue(valueNumber), nil
}

func expandOrdinalValue(valueNumber int) []string {
	if valueNumber < len(ordinalUnits) {
		return []string{ordinalUnits[valueNumber]}
	}
	tens := valueNumber / 10
	unit := valueNumber % 10
	if unit == 0 {
		return []string{ordinalTens[tens]}
	}
	return []string{cardinalTens[tens], ordinalUnits[unit]}
}

// ExpandPaul2013SlashDate expands an MM/DD/YYYY surface. The direct runtime
// capture is 01/02/2024 -> "January second twenty twenty four". Field widths
// and component bounds beyond that capture are conservative implementation
// assumptions; month-specific calendar validity has not been established.
func ExpandPaul2013SlashDate(value string) ([]string, error) {
	parts := strings.Split(value, "/")
	if len(parts) != 3 {
		return nil, errors.New("slash date must contain month, day, and year")
	}
	monthText, dayText, yearText := parts[0], parts[1], parts[2]
	if len(monthText) != 2 || len(dayText) != 2 || len(yearText) != 4 {
		return nil, errors.New("slash date must use MM/DD/YYYY field widths")
	}
	if err := validateDigits(monthText, "date month"); err != nil {
		return nil, err
	}
	if err := validateDigits(dayText, "date day"); err != nil {
		return nil, err
	}
	if err := validateDigits(yearText, "date year"); err != nil {
		return nil, err
	}
	month := int(monthText[0]-'0')*10 + int(monthText[1]-'0')
	day := int(dayText[0]-'0')*10 + int(dayText[1]-'0')
	months := [...]string{"", "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	if month < 1 || month >= len(months) {
		return nil, errors.New("date month is outside the supported 01 through 12 range")
	}
	if day < 1 || day > 31 {
		return nil, errors.New("date day is outside the supported 01 through 31 range")
	}
	yearWords, err := ExpandPaul2013UnsignedInteger(yearText)
	if err != nil {
		return nil, err
	}
	words := []string{months[month]}
	words = append(words, expandOrdinalValue(day)...)
	words = append(words, yearWords...)
	return words, nil
}

// ExpandPaul2013ClockTime spells a 12-hour H:MM clock token. The captured
// runtime example is 3:45 -> "three forty five"; ranges beyond that case are
// implemented by cardinal/minute formatting rules and are not DLL parity
// claims. A zero minute is rejected because its spoken form was not captured.
func ExpandPaul2013ClockTime(value string) ([]string, error) {
	hourText, minuteText, found := strings.Cut(value, ":")
	if !found || strings.Contains(minuteText, ":") {
		return nil, errors.New("clock time must contain exactly one colon")
	}
	if err := validateDigits(hourText, "clock hour"); err != nil {
		return nil, err
	}
	if err := validateDigits(minuteText, "clock minute"); err != nil {
		return nil, err
	}
	if len(hourText) > 2 || len(hourText) > 1 && hourText[0] == '0' {
		return nil, errors.New("clock hour must use one or two digits without a leading zero")
	}
	if len(minuteText) != 2 {
		return nil, errors.New("clock minute must use exactly two digits")
	}
	hour := int(hourText[0] - '0')
	if len(hourText) == 2 {
		hour = hour*10 + int(hourText[1]-'0')
	}
	minute := int(minuteText[0]-'0')*10 + int(minuteText[1]-'0')
	if hour < 1 || hour > 12 {
		return nil, errors.New("clock hour is outside the supported 1 through 12 range")
	}
	if minute < 1 || minute > 59 {
		return nil, errors.New("clock minute is outside the supported 1 through 59 range")
	}
	words := cardinalBelowThousand(hour)
	if minute < 10 {
		words = append(words, "oh", cardinalOnes[minute])
	} else {
		words = append(words, cardinalBelowThousand(minute)...)
	}
	return words, nil
}

func ordinalSuffix(value int) string {
	if value%100 >= 11 && value%100 <= 13 {
		return "th"
	}
	switch value % 10 {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}

func ungroupInteger(value string) (string, error) {
	if !strings.Contains(value, ",") {
		if err := validateDigits(value, "integer"); err != nil {
			return "", err
		}
		return value, nil
	}
	groups := strings.Split(value, ",")
	if len(groups[0]) < 1 || len(groups[0]) > 3 {
		return "", errors.New("grouping comma must follow one to three digits")
	}
	for index, group := range groups {
		if err := validateDigits(group, "integer group"); err != nil {
			return "", err
		}
		if index > 0 && len(group) != 3 {
			return "", errors.New("grouping comma must separate three-digit groups")
		}
	}
	return strings.Join(groups, ""), nil
}

func validateDigits(value, field string) error {
	if value == "" {
		return fmt.Errorf("%s contains no digits", field)
	}
	for position := 0; position < len(value); position++ {
		if value[position] < '0' || value[position] > '9' {
			return fmt.Errorf("%s contains non-digit byte at position %d", field, position)
		}
	}
	return nil
}

func expandFourDigitYear(value string) []string {
	firstTwo := int(value[0]-'0')*10 + int(value[1]-'0')
	hundredsDigit := value[2] - '0'
	lastTwo := int(value[2]-'0')*10 + int(value[3]-'0')
	if (firstTwo == 10 || firstTwo == 20) && lastTwo < 10 {
		return expandCardinalGroups(value)
	}
	words := cardinalBelowThousand(firstTwo)
	if lastTwo == 0 && value[1] != '0' {
		return append(words, "hundred")
	}
	if hundredsDigit == 0 && lastTwo > 0 && lastTwo < 10 {
		return append(words, "oh", cardinalOnes[lastTwo])
	}
	if lastTwo < 10 {
		return expandCardinalGroups(value)
	}
	return append(words, cardinalBelowThousand(lastTwo)...)
}

func cardinalBelowThousand(value int) []string {
	words := make([]string, 0, 5)
	if value >= 100 {
		words = append(words, cardinalOnes[value/100], "hundred")
		value %= 100
	}
	if value >= 20 {
		words = append(words, cardinalTens[value/10])
		value %= 10
	}
	if value > 0 {
		words = append(words, cardinalOnes[value])
	}
	return words
}
