package copier_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jinzhu/copier"
)

type methodConverterSource struct{}
type fieldConverterSource struct {
	Value string
	Extra bool
}

func (methodConverterSource) Value() string { return "source" }

func TestCopyMethodConverterError(t *testing.T) {
	sentinel := errors.New("conversion failed")
	for _, deepCopy := range []bool{false, true} {
		for _, fromMethod := range []bool{false, true} {
			for _, fail := range []bool{false, true} {
				name := fmt.Sprintf("deep=%t/method=%t/fail=%t", deepCopy, fromMethod, fail)
				t.Run(name, func(t *testing.T) {
					var source interface{} = fieldConverterSource{Value: "source"}
					if fromMethod {
						source = methodConverterSource{}
					}
					destination := struct{ Value string }{"unchanged"}
					calls := 0
					err := copier.CopyWithOption(&destination, source, copier.Option{
						DeepCopy: deepCopy,
						Converters: []copier.TypeConverter{{
							SrcType: "",
							DstType: "",
							Fn: func(src interface{}) (interface{}, error) {
								calls++
								if src != "source" {
									t.Errorf("converter input = %#v, want source", src)
								}
								if fail {
									return nil, sentinel
								}
								return "converted", nil
							},
						}},
					})
					if calls != 1 {
						t.Errorf("converter called %d times, want 1", calls)
					}
					want := "converted"
					if fail {
						want = "unchanged"
						if !errors.Is(err, sentinel) {
							t.Errorf("error = %v, want sentinel", err)
						}
					} else if err != nil {
						t.Errorf("unexpected error: %v", err)
					}
					if destination.Value != want {
						t.Errorf("destination = %q, want %q", destination.Value, want)
					}
				})
			}
		}
	}
}
