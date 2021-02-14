package internal

func ChunkSliceUints(slc []uint, size int) [][]uint {
	var slcLen = len(slc)
	var divided = make([][]uint, 0, (slcLen/size)+1)

	for i := 0; i < slcLen; i += size {
		end := i + size

		if end > slcLen {
			end = slcLen
		}

		divided = append(divided, slc[i:end])
	}

	return divided
}
