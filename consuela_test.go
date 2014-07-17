package consuela

import (
	"reflect"
	"testing"
)

type testCSVtoMapPair struct {
	filename    string
	mapresponce map[string]string
}

var csvToMapPair = []testCSVtoMapPair{
	{"./testingresources/list1.csv", map[string]string{"email1@fromlist1.csv": "", "email2@fromlist1.csv": "", "email3@fromlist1.csv": ""}},
	{"./testingresources/list2.csv", map[string]string{"email1@fromlist2.csv": "", "email2@fromlist2.csv": "", "email3@fromlist2.csv": ""}},
	{"./testingresources/list3.csv", map[string]string{"email1@fromlist3.csv": "", "email2@fromlist3.csv": "", "email3@fromlist3.csv": ""}},
}

func TestCSVMap(t *testing.T) {
	for _, pair := range csvToMapPair {
		csvmap := CsvToMap(pair.filename)
		if !reflect.DeepEqual(csvmap, pair.mapresponce) {
			t.Error(
				"For", pair.filename,
				"expected", pair.mapresponce,
				"got", csvmap,
			)
		}
	}
}

type testMerge struct {
	filenames []string
	mergedmap map[string]string
}

var mergetest = []testMerge{
	{[]string{"./testingresources/list1.csv", "./testingresources/list2.csv"}, map[string]string{"email1@fromlist1.csv": "", "email2@fromlist1.csv": "", "email3@fromlist1.csv": "", "email1@fromlist2.csv": "", "email2@fromlist2.csv": "", "email3@fromlist2.csv": ""}},
	{[]string{"./testingresources/list2.csv", "./testingresources/list3.csv"}, map[string]string{"email1@fromlist2.csv": "", "email2@fromlist2.csv": "", "email3@fromlist2.csv": "", "email1@fromlist3.csv": "", "email2@fromlist3.csv": "", "email3@fromlist3.csv": ""}},
	{[]string{"./testingresources/list1.csv", "./testingresources/list2.csv", "./testingresources/list3.csv"}, map[string]string{"email1@fromlist1.csv": "", "email2@fromlist1.csv": "", "email3@fromlist1.csv": "", "email1@fromlist2.csv": "", "email2@fromlist2.csv": "", "email3@fromlist2.csv": "", "email1@fromlist3.csv": "", "email2@fromlist3.csv": "", "email3@fromlist3.csv": ""}},
}

func TestMergeTest(t *testing.T) {
	for _, pair := range mergetest {
		tempmap := make([]map[string]string, len(mergetest), len(mergetest))
		for index := range pair.filenames {
			tempmap = append(tempmap, CsvToMap(pair.filenames[index]))
		}
		responce := mapMerge(tempmap)
		if !reflect.DeepEqual(responce, pair.mergedmap) {
			t.Error(
				"For", pair.filenames,
				"expected", pair.mergedmap,
				"got", responce,
			)
		}
	}
}

type testCompareListPair struct {
	filenames   []string
	mapresponse map[string]string
}

var comparelistPair = []testCompareListPair{
	{[]string{"./testingresources/list4.csv", "./testingresources/list1.csv"}, map[string]string{"email1@fromlist2.csv": "", "email3@fromlist2.csv": "", "email1@fromlist3.csv": "", "email3@fromlist3.csv": ""}},
	{[]string{"./testingresources/list4.csv", "./testingresources/list2.csv"}, map[string]string{"email1@fromlist1.csv": "", "email3@fromlist1.csv": "", "email1@fromlist3.csv": "", "email3@fromlist3.csv": ""}},
	{[]string{"./testingresources/list4.csv", "./testingresources/list3.csv"}, map[string]string{"email1@fromlist1.csv": "", "email3@fromlist1.csv": "", "email1@fromlist2.csv": "", "email3@fromlist2.csv": ""}},
}

func TestCompareList(t *testing.T) {
	for _, pair := range comparelistPair {
		csvmap := compareLists(CsvToMap(pair.filenames[0]), CsvToMap(pair.filenames[1]))
		if !reflect.DeepEqual(csvmap, pair.mapresponse) {
			t.Error(
				"For", pair.filenames[0], pair.filenames[1],
				"expected", pair.mapresponse,
				"got", csvmap,
			)
		}
	}
}
