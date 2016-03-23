package Image
import (
	"container/list"
	"image"
)

func Delaunay(point list.List, w int, l int) list.List{

	triangle_list := list.New()

	x := w
	y := l

	d := 40
	triangle_list.PushBack(NewTriangle(0 - d, 0 - d, x + d, 0 - d, x/2, y/2))
	triangle_list.PushBack(NewTriangle(x + d, y + d, x + d, 0 - d, x/2, y/2))
	triangle_list.PushBack(NewTriangle(0 - d, 0 - d, 0 - d, y + d, x/2, y/2))
	triangle_list.PushBack(NewTriangle(x + d, y + d, 0 - d, y + d, x/2, y/2))


	for e := point.Front(); e != nil; e = e.Next() {
		m_point := e.Value.(image.Point)
		int_rand_x := m_point.X
		int_rand_y := m_point.Y

		line_hash := make(map[Line]int)

		del_list := list.New()


		is_ignore := false
		for et := triangle_list.Front(); et != nil; et = et.Next() {
			tmp := et.Value.(*Triangle)
			if tmp.IsInSide(m_point) {
				is_ignore = true
				break
			}
		}
		if is_ignore {
			continue
		}

		for et := triangle_list.Front(); et != nil; et = et.Next() {
			tmp := et.Value.(*Triangle)
			if (tmp.IsInCircumcircle(image.Point{int_rand_x, int_rand_y})) {
				del_list.PushBack(et)
				line_hash[tmp.GetLine(0)]++
				line_hash[tmp.GetLine(1)]++
				line_hash[tmp.GetLine(2)]++
			}
		}

		for et := del_list.Front(); et != nil; et = et.Next() {
			tmp := et.Value.(*list.Element)
			triangle_list.Remove(tmp)
		}
		del_list.Init()

		for k, v := range line_hash {
			if v == 1 {
				triangle_list.PushBack(NewTriangle(int_rand_x, int_rand_y, k.Get(0).X, k.Get(0).Y, k.Get(1).X, k.Get(1).Y))
			}
		}

	}
	return *triangle_list
}