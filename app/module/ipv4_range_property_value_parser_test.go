package module_test

import (
	"github.com/odycenter/std-library/app/module"
	"reflect"
	"testing"
)

func TestIPv4RangePropertyValueParser_Parse(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "parseBlank",
			input: "",
			want:  []string{},
		},
		{
			name:  "parseCommaDelimited1",
			input: "cidr1, cidr2",
			want:  []string{"cidr1", "cidr2"},
		},
		{
			name:  "parseCommaDelimited2",
			input: " cidr1 ",
			want:  []string{"cidr1"},
		},
		{
			name:  "parseCommaDelimited3",
			input: "cidr1,cidr2 ",
			want:  []string{"cidr1", "cidr2"},
		},
		{
			name:  "parseSemicolonDelimited1",
			input: "name1: cidr1, cidr2; name2: cidr3,cidr4",
			want:  []string{"cidr1", "cidr2", "cidr3", "cidr4"},
		},
		{
			name:  "parseSemicolonDelimited2",
			input: "name1: cidr1; name2: cidr3; ",
			want:  []string{"cidr1", "cidr3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := module.NewIPv4RangePropertyValueParser(tt.input)
			got := parser.Parse()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
