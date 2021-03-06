package main
import ("fmt"; "log"; "strings"; "io/ioutil"; "time")

/** 
	User defined.
	
	@true prints various extra stuff out, but slows down the execution
	@false will be quick and quiet
*/
const debugMode bool = true

/**
 	Implementation of Basic Aho-Corasick algorithm (Prefix based).
	Searches for a set of strings (in 'patterns.txt') in text (in 'text.txt').
	Requires two files in the same folder as the algorithm:
	
	@file 'patterns.txt' containing the patterns to be searched for separated by single spaces
	@file 'text.txt' containing the text to be searched in
*/
func main() {
	patFile, err := ioutil.ReadFile("patterns.txt")
	if err != nil {
		log.Fatal(err)
	}
	textFile, err := ioutil.ReadFile("text.txt")
	if err != nil {
		log.Fatal(err)
	}
	patterns := strings.Split(string(patFile), " ")
	fmt.Printf("\nRunning: Basic Aho-Corasick algorithm.\n\n")
	if debugMode==true { 
		fmt.Printf("Searching for %d patterns/words:\n",len(patterns))
	}
	for i := 0; i < len(patterns); i++ {
		if (len(patterns[i]) > len(textFile)) {
			log.Fatal("There is a pattern that is longer than text! Pattern number:", i+1)
		}
		if debugMode==true { 
			fmt.Printf("%q ", patterns[i])
		}
	}
	if debugMode==true { 
		fmt.Printf("\n\nIn text (%d chars long): \n%q\n\n",len(textFile), textFile)
	}
	ahoCorasick(string(textFile), patterns)
}

/**
	Function performing the Basic Aho-Corasick alghoritm. 
	Finds and prints occurences of each pattern. 
	
	@param t text to be searched in
	@param p list of patterns to be serached for
*/  
func ahoCorasick(t string, p []string) {
	startTime := time.Now()
	occurences := make(map[int][]int)
	ac, f, s := buildAc(p)
	if debugMode==true {
		fmt.Printf("\n\nAC:\n\n")
	}
	current := 0
	for pos := 0; pos < len(t); pos++ {
		if debugMode==true {
			fmt.Printf("Position: %d, we read: %c", pos, t[pos])
        }
		for getTransition(current, t[pos], ac) == -1 && s[current] != -1 {
			current = s[current]
		}
		if getTransition(current, t[pos], ac) != -1 {
			current = getTransition(current, t[pos], ac)
			fmt.Printf(" (Continue) \n")
		} else {
			current = 0
			if debugMode==true {
				fmt.Printf(" (FAIL) \n")
			}
		}
		_, ok := f[current]
		if ok {
			for i := range f[current] {
				if p[f[current][i]] == getWord(pos-len(p[f[current][i]])+1, pos, t) { //check for word match
					if debugMode==true {
						fmt.Printf("Occurence at position %d, %q = %q\n", pos-len(p[f[current][i]])+1, p[f[current][i]], p[f[current][i]])
					}
					newOccurences := intArrayCapUp(occurences[f[current][i]])
					occurences[f[current][i]] = newOccurences
					occurences[f[current][i]][len(newOccurences)-1] = pos-len(p[f[current][i]])+1
				}
			}
		}
	}
	elapsed := time.Since(startTime)
	fmt.Printf("\n\nElapsed %f secs\n", elapsed.Seconds())
	for key, value := range occurences { //prints all occurences of each pattern (if there was at least one)
		fmt.Printf("\nThere were %d occurences for word: %q at positions: ",len(value), p[key])
		for i := range value {
			fmt.Printf("%d", value[i])
			if i != len(value) - 1 {
				fmt.Printf(", ")
			}
		}
		fmt.Printf(".")
	}
	return
}

/**
	Functions that builds Aho Corasick automaton.
*/
func buildAc(p []string) (acToReturn map[int]map[uint8]int, f map[int][]int, s []int) {
	acTrie, stateIsTerminal, f := constructTrie(p)
	s = make([]int, len(stateIsTerminal)) //supply function
	i := 0 //root of acTrie
	acToReturn = acTrie
	s[i] = -1
	if debugMode==true {
		fmt.Printf("\n\nAC construction: \n")
	}
	for current := 1; current < len(stateIsTerminal); current++ {
		o, parent := getParent(current, acTrie)
		down := s[parent]
		for stateExists(down, acToReturn) && getTransition(down, o, acToReturn) == -1 {
			down = s[down]
		}
		if stateExists(down, acToReturn) {
			s[current] = getTransition(down, o, acToReturn)
			if stateIsTerminal[s[current]] == true {
				stateIsTerminal[current] = true
				f[current] = arrayUnion(f[current], f[s[current]]) //F(Current) <- F(Current) union F(S(Current))
				if debugMode==true {
					fmt.Printf(" f[%d] set to: ", current)
					for i := range f[current] {
						fmt.Printf("%d\n", f[current][i])
					}
				}
			}
		} else {
			s[current] = i //initial state?
		}
	}
	if debugMode==true {
		fmt.Printf("\nsupply function: \n")
		for i:= range s {
			fmt.Printf("\ns[%d]=%d", i, s[i])
		}
		fmt.Printf("\n\n")
		for i,j := range f {
			fmt.Printf("f[%d]=", i)
			for k := range j {
				fmt.Printf("%d\n", j[k])
			}
		}
	}
	return acToReturn, f, s
}

/**
	Function that constructs Trie as an automaton for a set of reversed & trimmed strings.
	
	@return 'trie' built prefix tree
	@return 'stateIsTerminal' array of all states and boolean values of their terminality
	@return 'f' map with keys of pattern indexes and values - arrays of p[i] terminal states
*/
func constructTrie (p []string) (trie map[int]map[uint8]int, stateIsTerminal []bool, f map[int][]int) {
	trie = make(map[int]map[uint8]int)
	stateIsTerminal = make([]bool, 1)
	f = make(map[int][]int) 
	state := 1
	if debugMode==true {
		fmt.Printf("\n\nTrie construction: \n")
	}
	createNewState(0, trie)
	for i:=0; i<len(p); i++ {
		current := 0
		j := 0
		for j < len(p[i]) && getTransition(current, p[i][j], trie) != -1 {
			current = getTransition(current, p[i][j], trie)
			j++
		}
		for j < len(p[i]) {
			stateIsTerminal = boolArrayCapUp(stateIsTerminal)
			createNewState(state, trie)
			stateIsTerminal[state] = false
			createTransition(current, p[i][j], state, trie)
			current = state
			j++
			state++
		}
		if stateIsTerminal[current] {
			newArray := intArrayCapUp(f[current])
			newArray[len(newArray)-1] = i
			f[current] = newArray //F(Current) <- F(Current) union {i}
			if debugMode==true {
				fmt.Printf(" and %d", i)
			}
		} else {
			stateIsTerminal[current] = true
			f[current] = []int {i}  //F(Current) <- {i}
			if debugMode==true {
				fmt.Printf("\n%d is terminal for word number %d", current, i) 
			}
		}
	}
	return trie, stateIsTerminal, f
}

/**
	Returns 'true' if arry of int's 's' contains int 'e', 'false' otherwise.
	
	@author Mostafa http://stackoverflow.com/a/10485970
*/
func contains(s []int, e int) bool {
    for _, a := range s {
		if a == e {
			return true
		}
	}
    return false
}

/*******************          String functions          *******************/
/**
	Function that returns word found in text 't' at position range 'begin' to 'end'.
*/
func getWord(begin, end int, t string) string {
	for end >= len(t) {
		return ""
	}
	d := make([]uint8, end-begin+1)
	for j, i := 0, begin; i <= end; i, j = i+1, j+1 {
		d[j] = t[i]
	}
	return string(d)
}

/*******************   Array size allocation functions  *******************/
/**
	Dynamically increases an array size of int's by 1.
*/
func intArrayCapUp (old []int)(new []int) {
	new = make([]int, cap(old)+1)
	copy(new, old)  //copy(dst,src)
	old = new
	return new
}

/**
	Dynamically increases an array size of bool's by 1.
*/
func boolArrayCapUp (old []bool)(new []bool) {
	new = make([]bool, cap(old)+1)
	copy(new, old)
	old = new
	return new
}

/**
	Concats two arrays of int's into one.
*/
func arrayUnion (to, from []int) (concat []int) {
	concat = to
	for i := range(from) {
		if (!contains(concat, from[i])) {
			concat = intArrayCapUp(concat)
			concat[len(concat)-1] = from[i]
		}
	}
	return concat
}

/*******************          Automaton functions          *******************/
/**
	Function that finds the first previous state of a state and returns it. 
	Used for trie where there is only one parent.
	@param 'at' automaton
*/
func getParent(state int, at map[int]map[uint8]int) (uint8, int) {
	for beginState, transitions := range at {
		for c, endState := range transitions {
			if endState == state {
				return c, beginState
			}
		}
	}
	return 0, 0 //unreachable
}

/**
	Automaton function for creating a new state 'state'.
	@param 'at' automaton
*/
func createNewState(state int, at map[int]map[uint8]int) {
	at[state] = make(map[uint8]int)
	if debugMode==true {
		fmt.Printf("\ncreated state %d", state)
	}
}

/**
 	Creates a transition for function σ(state,letter) = end.
	@param 'at' automaton
*/
func createTransition(fromState int, overChar uint8, toState int, at map[int]map[uint8]int) {
	at[fromState][overChar]= toState
	if debugMode==true {
		fmt.Printf("\n    σ(%d,%c)=%d;",fromState,overChar,toState)
	}
}

/**
	Returns ending state for transition σ(fromState,overChar), '-1' if there is none.
	@param 'at' automaton
*/
func getTransition(fromState int, overChar uint8, at map[int]map[uint8]int)(toState int) {
	if (!stateExists(fromState, at)) {
		return -1
	}
	toState, ok := at[fromState][overChar]
	if (ok == false) {
		return -1	
	}
	return toState
}

/**
	Checks if state 'state' exists. Returns 'true' if it does, 'false' otherwise.
	@param 'at' automaton
*/
func stateExists(state int, at map[int]map[uint8]int)bool {
	_, ok := at[state]
	if (!ok || state == -1 || at[state] == nil) {
		return false
	}
	return true
}