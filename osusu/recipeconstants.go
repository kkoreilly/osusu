package osusu

// BaseCuisines are the base cuisines that are assigned to recipes and can be used in meals
var BaseCuisines = []string{"African", "American", "Asian", "British", "Chinese", "European", "French", "Greek", "Indian", "Italian", "Japanese", "Jewish", "Korean", "Latin American", "Mexican", "Middle Eastern", "Thai"}

// CuisineToCuisineMap has a list of all cuisines used and their aliases
var CuisineToCuisineMap = map[string][]string{
	"African":        {"African", "East African", "North African", "South African", "West African", "Algerian", "Egyptian", "Ethiopian", "Moroccan", "Tunisian"},
	"American":       {"American", "New England", "North American", "Tex Mex", "Tex-Mex", "U.S.", "Pennsylvania Dutch", "Native American", "Southern", "Southwestern", "Amish", "Hawaiian", "Cajun", "Creole", "Canadian", "French Canadian"},
	"Asian":          {"Asian", "Asian Inspired", "East And Southeast Asian", "South And Central Asian", "Bangladeshi", "Filipino", "Indonesian", "Malaysian", "Sri Lankan", "Armenian", "Vietnamese"},
	"British":        {"British", "English", "Uk And Ireland", "Scottish", "Welsh", "Irish", "Oceanic", "Australian", "New Zealand"},
	"Chinese":        {"Chinese", "Chinese Inspired", "Sichuan"},
	"European":       {"European", "Eastern European", "Russian", "Ukrainian", "Finnish", "Norwegian", "Scandinavian", "Swedish", "Austrian", "Belgian", "Czech", "Danish", "Dutch", "Hungarian", "Polish", "Portuguese", "Swiss", "Turkish", "German", "Spanish"},
	"French":         {"French", "French Inspired"},
	"Greek":          {"Greek", "Greek Inspired", "Mediterranean Inspired"},
	"Indian":         {"Indian", "Indian Inspired"},
	"Italian":        {"Italian", "Italian Inspired", "Sicilian"},
	"Japanese":       {"Japanese", "Japanese Inspired"},
	"Jewish":         {"Jewish", "Kosher", "Israeli"},
	"Korean":         {"Korean", "Korean Inspired"},
	"Latin American": {"Latin American", "Latin", "South American", "Argentine", "Venezuelan", "Brazilian", "Caribbean", "Chilean", "Colombian", "Cuban", "Puerto Rican", "Salvadoran", "Peruvian", "Jamaican"},
	"Mexican":        {"Mexican", "Mexican Inspired"},
	"Middle Eastern": {"Middle Eastern", "Middle Eastern Inspired", "Afghan", "Lebanese", "Pakistani", "Syrian", "Persian"},
	"Thai":           {"Thai", "Thai Inspired"},
}

// IgnoredCuisines are the allrecipes cuisines that we ignore
var IgnoredCuisines = []string{
	"Fusion", "Inspired", "World", "Copycat", "Authentic",
}

// AllCuisines are all of the cuisines possible
var AllCuisines = []string{"African", "American", "Anglo-Indian", "Arabian", "Argentine", "Armenian", "Australian", "Austrian", "Azeri",
	"Balkan", "Bangladeshi", "Barbeque", "Basque", "Belgian", "Bengali", "Bhutanese", "Bolivian", "Brazilian", "British",
	"Bruneian", "Bulgarian", "Burmese", "Cambodian", "Cantonese", "Cape Malay", "Central Asian", "Cherokee", "Chilean",
	"Chinese", "Colombian", "Cornish", "Costa Rican", "Croatian", "Cuban", "Cypriot", "Czech", "Danish", "Djiboutian",
	"Dominican", "Dutch", "East African", "Eastern European", "Ecuadorian", "Egyptian", "Eritrean", "Estonian",
	"Ethiopian", "Faroe Islands", "Filipino", "Finnish", "French", "Galician", "Gambian", "Georgian", "German",
	"Ghanaian", "Greek", "Grenadian", "Guatemalan", "Guinea-Bissauan", "Guyanese", "Haitian", "Hawaiian", "Herzegovinian",
	"Hungarian", "Icelandic", "Indian", "Indonesian", "Iranian", "Iraqi", "Irish", "Israeli", "Italian",
	"Jamaican", "Japanese", "Jordanian", "Kazakh", "Kenyan", "Khmer", "Korean", "Kosovan", "Kuwaiti",
	"Kyrgyz", "Laotian", "Latin American", "Latvian", "Lebanese", "Lithuanian", "Luxembourgish", "Macedonian",
	"Malagasy", "Malaysian", "Maldivian", "Maltese", "Marshallese", "Mauritanian", "Mauritian", "Mexican",
	"Micronesian", "Middle Eastern", "Mongolian", "Moroccan", "Mozambican", "Myanmar", "Namibian", "Nepalese",
	"New Zealand", "Nicaraguan", "Nigerian", "North African", "North American", "Norwegian", "Omani", "Pakistani",
	"Palauan", "Palestinian", "Panamanian", "Papua New Guinean", "Paraguayan", "Peruvian", "Philippine",
	"Polish", "Portuguese", "Qatari", "Romanian", "Russian", "Rwandan", "Saint Lucian", "Salvadoran", "Samoa",
	"Samoan", "Sanmarinese", "Sao Tome and Principe", "Saudi Arabian", "Scottish", "Senegalese", "Serbian",
	"Seychellois", "Sierra Leonean", "Singaporean", "Slovak", "Slovenian", "Solomon Islander", "Somali",
	"South African", "South American", "South Korean", "Spanish", "Sri Lankan", "Sudanese", "Surinamese",
	"Swazi", "Swedish", "Swiss", "Syrian", "Taiwanese", "Tajikistani", "Tanzanian", "Thai", "Tibetan",
	"Tonga", "Trinidad and Tobago", "Tunisian", "Turkish", "Turkmen", "Tuvaluan", "Ugandan", "Ukrainian",
	"Uruguayan", "Uzbek", "Vietnamese", "Welsh", "West African", "Western European", "Yemeni", "Zambian",
	"Zimbabwean"}
