package gogo

func (b *Board) Score() float64 {
    worms := b.MakeWorms()
    ret := 0.0
    for _, worm := range worms {
        if worm.Color == GRAY {
            if worm.BorderColor == BLACK {
                ret += float64(worm.Points.Size())
            } else if worm.BorderColor == WHITE {
                ret -= float64(worm.Points.Size())
            }
        } else if worm.Color == BLACK {
            ret += float64(worm.Points.Size())
        } else if worm.Color == WHITE {
            ret -= float64(worm.Points.Size())
        }
    }
    return ret - b.komi
}
