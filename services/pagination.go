package services

import "math"

type PaginationService struct {
	minPageSize int
	maxPageSize int
	pageSize    int
	pageNumber  int
	numItems    int
}

// NewPaginationService returns a new pagination service with pagination
// defaults
func NewPaginationService() PaginationService {
	p := PaginationService{}
	p.SetMaxPageSize(25)
	p.SetMinPageSize(10)
	p.SetPageSize(15)

	return p
}

// SetMinPageSize sets the min page size
func (p *PaginationService) SetMinPageSize(newMin int) {
	p.minPageSize = newMin
}

// SetMaxPageSize sets the max page size
func (p *PaginationService) SetMaxPageSize(newMax int) {
	p.maxPageSize = newMax
}

// SetPageSize sets the page size while making sure it is within the limits of
// the max and min page size
func (p *PaginationService) SetPageSize(size int) {
	if size > p.maxPageSize {
		size = p.maxPageSize
	} else if size < p.minPageSize {
		size = p.minPageSize
	}

	p.pageSize = size
}

func (p *PaginationService) SetPageNumber(num int) {
	if num <= 0 {
		num = 1
	}

	p.pageNumber = num
}

func (p *PaginationService) SetNumItems(num int) {
	p.numItems = num
}

func (p *PaginationService) GetNumItems() int {
	return p.numItems
}

func (p *PaginationService) GetStartNumber() int {
	// make sure the page number is not 0
	p.SetPageNumber(p.pageNumber)

	return p.pageNumber*p.pageSize - p.pageSize
}

func (p *PaginationService) GetFirst() int {
	return 1
}

func (p *PaginationService) GetLast() int {
	return int(math.Ceil(float64(p.numItems) / float64(p.pageSize)))
}

func (p *PaginationService) GetPrev() int {
	if p.pageNumber <= 1 {
		return 0
	} else {
		return p.pageNumber - 1
	}
}

func (p *PaginationService) GetNext() int {
	if p.pageNumber >= p.GetLast() {
		return p.pageNumber
	} else {
		return p.pageNumber + 1
	}
}

func (p *PaginationService) GetCurrent() int {
	return p.pageNumber
}

func (p *PaginationService) GetSize() int {
	return p.pageSize
}
