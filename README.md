# png2c

A simple tool for converting PNG images to a C definition.

## Usage

Download and install the library and tools using:

```
go get github.com/silkeh/png2c
```

To convert an image, use it as follows:

```
png2c -file pixel.png -var pixel -mode 565 -brief Pixel
```

Which will result in something like:

```c
/**
 * @brief   Pixel (1x1)
 */
const uint16_t pixel[1][1] = {
    { 0xe005 }
};
```

Note that no compensation for the alpha channel is done when the alpha channel
is not represented in the output format.
