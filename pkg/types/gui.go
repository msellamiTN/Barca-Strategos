package types

// InterfaceComplexity defines UI complexity levels.
type InterfaceComplexity string

const (
	InterfaceComplexitySimple  InterfaceComplexity = "simple"
	InterfaceComplexityMedium InterfaceComplexity = "medium"
	InterfaceComplexityComplex InterfaceComplexity = "complex"
)

// LayoutType defines UI layout options.
type LayoutType string

const (
	LayoutTypeGrid    LayoutType = "grid"
	LayoutTypeList    LayoutType = "list"
	LayoutTypeCard    LayoutType = "card"
	LayoutTypeTable   LayoutType = "table"
)

// ThemeType defines UI themes.
type ThemeType string

const (
	ThemeTypeLight ThemeType = "light"
	ThemeTypeDark  ThemeType = "dark"
	ThemeTypeAuto  ThemeType = "auto"
)
