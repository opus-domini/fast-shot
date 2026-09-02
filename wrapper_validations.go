package fastshot

// Compile-time check that DefaultValidations implements ValidationsWrapper.
var _ ValidationsWrapper = (*DefaultValidations)(nil)

// DefaultValidations implements ValidationsWrapper interface and provides a default HTTP validations.
type DefaultValidations struct {
	validations []error
}

// Unwrap will return the underlying validations
func (c *DefaultValidations) Unwrap() []error {
	return c.validations
}

// Get will return the validation at the specified index
func (c *DefaultValidations) Get(index int) error {
	if index < 0 || index >= len(c.validations) {
		return nil
	}
	return c.validations[index]
}

// IsEmpty will return the underlying validations
func (c *DefaultValidations) IsEmpty() bool {
	return len(c.validations) == 0
}

// Count will return the number of validations
func (c *DefaultValidations) Count() int {
	return len(c.validations)
}

// Add will append a new validation to the underlying validations
func (c *DefaultValidations) Add(err error) {
	c.validations = append(c.validations, err)
}

// newDefaultValidations initializes a new DefaultValidations with a given validations.
func newDefaultValidations(validations []error) *DefaultValidations {
	if validations == nil {
		validations = make([]error, 0)
	}
	return &DefaultValidations{
		validations: validations,
	}
}
