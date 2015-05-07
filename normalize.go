package vcfgo

func leftalign(pos int, ref []byte, alt []byte, seq []byte) (int, []byte, []byte, error) {
	/* actually this isn't necessary
	if !bytes.HasSuffix(seq, ref) {
		return 0, ref, alt, errors.New("leftalign: sequence should end with ref")
	}
	*/
	subseq := seq[:len(seq)-len(ref)]

	quit := false
	j := len(subseq)

	for j > 0 && !quit {
		quit = true
		if ref[len(ref)-1] == alt[len(alt)-1] {
			ref = ref[:len(ref)-1]
			alt = alt[:len(alt)-1]
			quit = false
		}
		if len(ref) == 0 || len(alt) == 0 {
			j--
			// use j + 1 to get a slice.
			ref = append(subseq[j:j+1], ref...)
			alt = append(subseq[j:j+1], alt...)
			quit = false

		}

	}
	return pos - (len(subseq) - j), ref, alt, nil
}

func lefttrim(pos int, ref []byte, alt []byte) (int, []byte, []byte, error) {
	if 1 == len(ref) && len(alt) == 1 {
		return pos, ref, alt, nil
	}

	n := 0
	minLen := len(ref)
	if len(alt) < minLen {
		minLen = len(alt)
	}

	for n+1 < minLen && alt[n] == ref[n] {
		n++
	}

	alt, ref = alt[n:], ref[n:]
	pos += n
	return pos, ref, alt, nil
}
