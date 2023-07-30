# Notes

## Benchmarking Lines.Render

1. image.NewNRGBA(...)                152454   ns/op   0.15 ms/op
2. Fill with transparent background   2908475  ns/op   2.91 ms/op
3. Draw blank waveform                10380673 ns/op  10.38 ms/op
4. Antialias blank                    38450586 ns/op  38.45 ms/op
5. Summing audio samples              40439157 ns/op  40.44 ms/op
6. Colouring pixels                   47014037 ns/op  47.01 ms/op

## SIMD

1. https://stackoverflow.com/questions/32083301/equivalent-simd-instruction-for-multiplying-specific-array-elements#32093588