package main

import "strconv"

// port is a convience type to allow the acceptance of uint16 as a flag value
type port uint16

func (p port) String() string {
	return strconv.Itoa(int(p))
}

func (p *port) Set(s string) error {
	var (
		i   uint64
		err error
	)

	// ParseUint handles bounds checking
	if i, err = strconv.ParseUint(s, 10, 16); err != nil {
		return err
	}

	*p = port(i)
	return nil

}

func (p port) Get() interface{} {
	return p
}
