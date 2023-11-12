package osusu

// AllCategories are all of the possible meal categories
var AllCategories = []string{"Breakfast", "Brunch", "Lunch", "Dinner", "Dessert", "Snack", "Appetizer", "Side", "Drink", "Ingredient"}

// AllSources are all of the possible meal sources
var AllSources = []string{"Cooking", "Dine-In", "Takeout", "Delivery"}

// AllDietaryRestrictions are all of the default possible dietary restrictions. Users can add specific ingredients that they can't have separately if needed.
var AllDietaryRestrictions = []string{"Vegetarian", "Vegan", "Gluten-Free", "Lactose-Free", "Nut-Free", "Dairy-Free", "Kosher", "Halal"}

// DietaryRestrictionsIngredientsMap are the ingredients disallowed for each dietary restriction in AllDietaryRestrictions
var DietaryRestrictionsIngredientsMap = map[string][]string{
	"Vegetarian": {
		"meat", "pork", "beef", "veal", "lamb", "sausage", "sausages", "bacon", "goat", "venison", "rabbit", "pepperoni", "steak", "steaks", "ribs", "ham", "chuck", "dogs", "prosciutto", "brisket",
		"poultry", "chicken", "turkey", "duck", "emu", "goose",
		"fish", "prawns", "mussels", "oysters", "clams", "salmon", "tuna", "scallops", "halibut", "lobster", "crab", "flounder", "trout", "shrimp", "catfish", "crawfish", "octopus",
	},
}

// CategoryToCategoryMap has a list of all categories and their aliases
var CategoryToCategoryMap = map[string][]string{
	"Dinner":     {"Dinner", "Entree", "Pasta"},
	"Drink":      {"Drink", "Beverage", "Cocktail", "Coffee"},
	"Dessert":    {"Dessert", "Cake", "Candy", "Pie"},
	"Lunch":      {"Lunch", "Sandwich"},
	"Ingredient": {"Ingredient", "Bread", "Condiment", "Jam / Jelly", "Sauce", "Spice Mix"},
	"Appetizer":  {"Appetizer", "Salad", "Soup"},
	"Side":       {"Side", "Side Dish"},
	"Breakfast":  {"Breakfast"},
	"Snack":      {"Snack"},
	"Brunch":     {"Brunch"},
}

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

// IgnoredWords are all of the non-meaningful words that are excluded from the GetWords function
var IgnoredWords = []string{
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "30", "35", "40", "45", "50", "55", "60", "65", "70", "75", "80", "85", "90", "95", "100", "125", "150", "175", "200", "225", "250", "1/2", "1/3", "1/6",

	"cup", "cups", "teaspoon", "teaspoons", "tablespoon", "tablespoons", "ounce", "ounces", "chopped", "ground", "taste", "fresh", "sliced", "minced", "diced", "pound", "pounds", "cut", "dried", "shredded", "peeled", "divided", "pinch", "drained", "large", "medium", "small", "grated", "softened", "package", "finely", "whole", "slices", "crushed", "pieces", "melted", "freshly", "mix", "inch", "thinly", "beaten", "cubed", "-", "cooked", "spray", "hot", "halved", "sweet", "salt", "sugar", "pepper", "white", "oil", "sauce", "black", "flour", "cream", "powder", "water", "milk", "all-purpose", "recipe", "red", "vanilla", "juice", "green", "easy", "baking", "extract", "brown", "lemon", "make",

	"such", "needed", "This", "A", "The", "These",

	"a", "about", "above", "after", "again", "against", "ago", "ahead", "all", "almost", "along", "already", "also", "although", "always", "am", "among", "an", "and", "any", "are", "aren't", "around", "as", "at", "away",
	"backward", "backwards", "be", "because", "before", "behind", "below", "beneath", "beside", "between", "both", "but", "by",
	"can", "cannot", "can't", "cause", "cos", "could", "couldn't",
	"d", "despite", "did", "didn't", "do", "does", "doesn't", "don't", "down", "during",
	"each", "either", "even", "ever", "every", "except",
	"for", "forward", "from",
	"had", "hadn't", "has", "hasn't", "have", "haven't", "he", "her", "here", "hers", "herself", "him", "himself", "his", "how", "however", "I",
	"if", "in", "inside", "inspite", "instead", "into", "is", "isn't", "it", "its", "itself",
	"just",
	"ll", "least", "less", "like",
	"m", "many", "may", "mayn't", "me", "might", "mightn't", "mine", "more", "most", "much", "must", "mustn't", "my", "myself",
	"near", "need", "needn't", "needs", "neither", "never", "no", "none", "nor", "not", "now",
	"of", "off", "often", "on", "once", "only", "onto", "or", "ought", "oughtn't", "our", "ours", "ourselves", "out", "outside", "over",
	"past", "perhaps",
	"quite",
	"re", "rather",
	"s", "seldom", "several", "shall", "shan't", "she", "should", "shouldn't", "since", "so", "some", "sometimes", "soon",
	"than", "that", "the", "their", "theirs", "them", "themselves", "then", "there", "therefore", "these", "they", "this", "those", "though", "through", "thus", "till", "to", "together", "too", "towards",
	"under", "unless", "until", "up", "upon", "us", "used", "usedn't", "usen't", "usually",
	"ve", "very",
	"was", "wasn't", "we", "well", "were", "weren't", "what", "when", "where", "whether", "which", "while", "who", "whom", "whose", "why", "will", "with", "without", "won't", "would", "wouldn't",
	"yet", "you", "your", "yours", "yourself", "yourselves",
}

// IgnoredWordsMap is a map version of IgnoredWords with the key as the word and the value as true. It is initialized by InitRecipeConstants.
var IgnoredWordsMap = map[string]bool{}

// WordSeparators are the characters that can separate words
var WordSeparators = []rune{' ', ',', '.', '(', ')', '+', '–', '—'}

// WordSeparatorsMap is a map version of WordSeparators with the key as the separator and the value as true. It is initialized by InitRecipeConstants.
var WordSeparatorsMap = map[rune]bool{}

// InitRecipeConstants initializes IgnoredWordsMap and WordSeparatorsMap. It should be called once on program start.
func InitRecipeConstants() {
	for _, word := range IgnoredWords {
		IgnoredWordsMap[word] = true
	}
	for _, separator := range WordSeparators {
		WordSeparatorsMap[separator] = true
	}
}
