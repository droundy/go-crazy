// This is a nice example...

package main

type Vec []float64

func (a Vec) .- (b Vec) Vec {
	return Vec{ a[0]-b[0], a[1]-b[1], a[2]-b[2] }
}

func (a Vec) *. (b float64) Vec {
	return Vec{ a[0]*b, a[1]*b, a[2]*b }
}

func (a Vec) .+= (b Vec) {
	// This isn't so efficient, but demonstrates how these method
	// operators are used.
	a = a .+ b
}

func test(a, b, c Vec) Vec {
	return a .+ (b .+ c) .- a
}

func (a Vec) .+ (b Vec) Vec {
	return Vec{ a[0]+b[0], a[1]+b[1], a[2]+b[2] }
}

func main() {
	x := Vec{1,2,3}
	y := Vec{3,2,1}
	if (x .- y)[0] != -2 {
		panic("bug!")
	}
	if (x .- y)[1] != 0 {
		panic("bug!")
	}
	if (x .- y)[2] != 2 {
		panic("bug!")
	}
	if (2 *. x)[0] != 2 {
		panic("bug")
	}
}
