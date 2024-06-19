For a given JSON Data stream containing words which can have more than one occurances and can have Upper or Lower case characters,
convert it into a frequency map output in JSON format with words sorted lexicographically. 
Also, convert all the characters of the words to lowercase and sort 

Sample input: {"text": "cat Mat bat Rat Cat cat Bat"}
Sample Output: [{"w":"bat","c":2},{"w":"cat","c":3},{"w":"mat","c":1},{"w":"rat","c":1}]

Logic/Algorithm:
 1. Unmarshall the Input Json byte slice into a map to extract the string of words
    (Could use a struct of 2 strings as well since we know the type of key and value in this case, but map will also work)
 2. Use the value of unmarshalled output, which is the string containing multiple words, as input for further processing
    Convert this string to lowercase string
    Tokenize the string to get the list of words
    Store these words in a map with Key = word, Value = Frequency/count
    Create a new entry in the map when the key doesnt exits and Increment the frequency if the entry already exists
 3. In order to maintain lexicographic order for the output, we need to create a temporary slice of words
    because iterating over a map produces randomly ordered output
    Sort the temporary slice containing all the words fetched by iterating over map key values
    Iterate over sorted temporary word slice
    Fetch the frequency of that word from map using word as the key
    Create a slice of Json-data struture elements and instantiate using values of word and frequency for w and c respectively
    Marshall the structure to get the final JSON output
