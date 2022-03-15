package todo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	testCases := []struct {
		testName  string
		taskNames []string
	}{
		{testName: "Add 1 Item", taskNames: []string{"item1"}},
		{testName: "Add 2 of the Same Items", taskNames: []string{"item1", "item1"}},
		{testName: "Add 3 Distinct Items", taskNames: []string{"item1", "item2", "item3"}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			list := List{}
			var expected List

			for _, taskName := range testCase.taskNames {
				list.Add(taskName)
				expected = append(expected, item{
					Task:        taskName,
					Done:        false,
					CreatedAt:   time.Now(),
					CompletedAt: time.Time{},
				})
			}

			assertLists(t, list, expected, time.Millisecond, time.Millisecond)
		})
	}
}

func TestComplete(t *testing.T) {
	testCases := []struct {
		testName      string
		itemNums      []int
		errExpected   bool
		completeTwice bool
	}{
		{testName: "Invalid Item Number (Less)", itemNums: []int{0}, errExpected: true, completeTwice: false},
		{testName: "Invalid Item Number (Greater)", itemNums: []int{20}, errExpected: true, completeTwice: false},
		{testName: "Complete 1 Item", itemNums: []int{1}, errExpected: false, completeTwice: false},
		{testName: "Complete Same Item Twice", itemNums: []int{1, 1}, errExpected: true, completeTwice: true},
		{testName: "Complete All Items", itemNums: []int{1, 2, 3}, errExpected: false, completeTwice: false},
	}

	setupList := func() (list List) {
		list.Add("item1")
		list.Add("item2")
		list.Add("item3")

		return list
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			list := setupList()
			expected := setupList()

			for i, itemNum := range testCase.itemNums {
				err := list.Complete(itemNum)

				if testCase.errExpected && !testCase.completeTwice {
					assert.EqualError(t, err, ErrItemNotFound.Errorf(itemNum).Error())
				} else {
					if testCase.completeTwice && i > 0 {
						assert.EqualError(t, err, ErrItemAlreadyCompleted.Errorf(itemNum).Error())
					} else {
						assert.Nil(t, err)
						currItem := &expected[itemNum-1]
						currItem.Done = true
						currItem.CompletedAt = time.Now()
					}

					assertLists(t, list, expected, 500*time.Microsecond, 500*time.Microsecond)
				}

				time.Sleep(time.Millisecond)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	testCases := []struct {
		testName    string
		itemNums    []int
		errExpected bool
	}{
		{testName: "Invalid Item Num (Less)", itemNums: []int{0}, errExpected: true},
		{testName: "Invalid Item Num (Greater)", itemNums: []int{20}, errExpected: true},
		{testName: "Delete Only One Item", itemNums: []int{3}, errExpected: false},
		{testName: "Delete Item 1 Twice", itemNums: []int{1, 1}, errExpected: false},
		{testName: "Delete All Items", itemNums: []int{1, 2, 1}, errExpected: false},
		{testName: "Delete All Items (Reverse)", itemNums: []int{3, 2, 1}, errExpected: false},
	}

	setupList := func() (list List) {
		list.Add("item1")
		list.Add("item2")
		list.Add("item3")

		return list
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			list := setupList()
			prev := setupList()

			for _, itemNum := range testCase.itemNums {
				curr := List{}

				err := list.Delete(itemNum)

				if testCase.errExpected {
					assert.EqualError(t, err, ErrItemNotFound.Errorf(itemNum).Error())
				} else {
					assert.Nil(t, err)
					for i, item := range prev {
						if i != itemNum-1 {
							curr.Add(item.Task)
						}
					}

					assertLists(t, list, curr, time.Millisecond, time.Millisecond)
					prev = curr
					curr = nil
				}
			}
		})
	}
}

func assertLists(t testing.TB, got, expected List, createdTimeGap, completedTimeGap time.Duration) {
	t.Helper()

	assert.Equal(t, len(got), len(expected))

	for i := range got {
		assert.Equal(t, expected[i].Task, got[i].Task)
		assert.Equal(t, expected[i].Done, got[i].Done)
		assert.LessOrEqual(t, absTimeDuration(t, expected[i].CreatedAt, got[i].CreatedAt), createdTimeGap)
		assert.LessOrEqual(t, absTimeDuration(t, expected[i].CompletedAt, got[i].CompletedAt), completedTimeGap)
	}
}

func absTimeDuration(t testing.TB, timeA, timeB time.Time) time.Duration {
	t.Helper()

	if timeA.After(timeB) {
		return timeA.Sub(timeB)
	}

	return timeB.Sub(timeA)
}
