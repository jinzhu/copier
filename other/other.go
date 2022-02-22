package other

type TestEnum int64

const (
	Hello TestEnum = 1
	Hey   TestEnum = 2
)

type SecondLayer struct {
	Test   *TestEnum
	Others *string
}
