package darts

import (
	"bufio"
	"os"
	"sort"
	"strings"
	"testing"
)

var benchData = []string{
	"簳面杖", "簸荡", "簿历", "簿录", "簿据", "籧篨戚施", "米价", "米克", "米克杰格",
	"米克森", "米克诺斯", "米利托", "米制", "米卤蛋", "米厘米突", "米德尔伯里",
	"米格式战斗机", "米纳谷", "米罗的维纳斯雕像", "米苏里", "米苏里州", "米虫",
	"米蛀虫", "米谷", "米里", "米雅托维奇", "米雕", "米面", "类似于", "类别",
	"类别的团体", "类同", "类同法", "类球面", "粉丝团", "粉丝谷", "粉团儿",
	"粉团儿似的", "粉彩", "粉拳绣腿", "粉板", "粉砂岩", "粉签子", "粉红色系",
	"粉面", "粉面朱脣", "粉面油头", "粉饰门面", "粒变岩", "粗制", "粗制品",
	"粗制滥造", "粗卤", "粗布", "粗布条", "粗恶", "粗枝大叶",
	"粗毛布", "粗管面", "粗纤维", "粗衣恶食", "粗面", "粗面岩",
}

func TestMain(m *testing.M) {
	sort.Strings(benchData)
	m.Run()
}

func TestThrough(t *testing.T) {
	builder := DoubleArrayBuilder{}
	keyset := []string{"印度", "印度尼西亚", "印加帝国", "瑞士", "瑞典", "巴基斯坦", "巴勒斯坦", "以色列", "巴比伦", "土耳其"}
	sort.Strings(keyset)
	t.Log(keyset)

	builder.Build(keyset)

	// Check keys that exist
	for i, key := range keyset {
		result, matched := builder.ExactMatchSearch(key)
		if !matched {
			t.Errorf("invalid result, not matched: %s", key)
			return
		}

		if result != i {
			t.Errorf("invalid result, expected: %d, got: %d for key: %s", i, result, key)
			return
		}
	}

	// Check keys not exist
	notExists := []string{"印", "印度尼", "日本", "tokyo", "Turkey", "吐", "以色烈"}
	for _, key := range notExists {
		_, matched := builder.ExactMatchSearch(key)
		if matched {
			t.Errorf("invalid result, should not match: %s", key)
			return
		}
	}

	// Check Prefix search
	values := builder.CommonPrefixSearch("印度尼西亚啊")
	if len(values) != 2 || values[0] != 2 || values[1] != 3 {
		t.Errorf("common prefix searche result error, expected: [1, 2], got: %v", values)
	}

	byts := builder.ToBytes()

	da := DoubleArray{}
	err := da.FromBytes(byts)
	if err != nil {
		t.Errorf("can't build double array from bytes")
		return
	}

	// Check keys that exist
	for i, key := range keyset {
		result, matched := da.ExactMatchSearch(key)
		if !matched {
			t.Errorf("invalid result, not matched: %s", key)
			return
		}

		if result != i {
			t.Errorf("invalid result, expected: %d, got: %d for key: %s", i, result, key)
			return
		}
	}

	// Check keys not exist
	for _, key := range notExists {
		_, matched := da.ExactMatchSearch(key)
		if matched {
			t.Errorf("invalid result, should not match: %s", key)
			return
		}
	}

	values = da.CommonPrefixSearch("印度尼西亚啊")
	if len(values) != 2 || values[0] != 2 || values[1] != 3 {
		t.Errorf("common prefix searche result error, expected: [1, 2], got: %v", values)
	}
}

func BenchmarkDoubleArray(b *testing.B) {
	builder := DoubleArrayBuilder{}
	builder.Build(benchData)
	size := len(benchData)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		idx := n % size
		result, matched := builder.ExactMatchSearch(benchData[idx])
		if !matched {
			b.Errorf("not matched: %s", benchData[idx])
			return
		}
		if result != idx {
			b.Errorf("error result, expected: %d, got: %d", idx, result)
			return
		}
	}
}

func BenchmarkMap(b *testing.B) {
	m := make(map[string]int)
	for i, key := range benchData {
		m[key] = i
	}
	size := len(benchData)

	for n := 0; n < b.N; n++ {
		idx := n % size
		result, matched := m[benchData[idx]]
		if !matched {
			b.Errorf("not matched: %s", benchData[idx])
			return
		}
		if result != idx {
			b.Errorf("error result, expected: %d, got: %d", idx, result)
			return
		}
	}
}

func keysFromFile() (keys []string, err error) {
	file, err := os.Open("testdata/STPhrases.txt")
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	keys = make([]string, 0, 10000)
	for scanner.Scan() {
		text := scanner.Text()
		text = strings.TrimSpace(text)
		kv := strings.Split(text, "\t")
		if len(kv) != 2 {
			continue
		}
		keys = append(keys, kv[0])
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	sort.Strings(keys)
	return keys, nil
}

func BenchmarkBigDoubleArray(b *testing.B) {
	benchData, err := keysFromFile()
	if err != nil {
		b.Error(err)
		return
	}

	builder := DoubleArrayBuilder{}
	builder.Build(benchData)
	size := len(benchData)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		idx := n % size
		result, matched := builder.ExactMatchSearch(benchData[idx])
		if !matched {
			b.Errorf("not matched: %s", benchData[idx])
			return
		}
		if result != idx {
			b.Errorf("error result, expected: %d, got: %d", idx, result)
			return
		}
	}
}

func BenchmarkBigMap(b *testing.B) {
	benchData, err := keysFromFile()
	if err != nil {
		b.Error(err)
		return
	}

	m := make(map[string]int)
	for i, key := range benchData {
		m[key] = i
	}
	size := len(benchData)

	for n := 0; n < b.N; n++ {
		idx := n % size
		result, matched := m[benchData[idx]]
		if !matched {
			b.Errorf("not matched: %s", benchData[idx])
			return
		}
		if result != idx {
			b.Errorf("error result, expected: %d, got: %d", idx, result)
			return
		}
	}
}
