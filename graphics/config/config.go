package config

type PanningOptions struct {
    PanXSize uint16
    PanYSize uint16

    StartingX uint16
    StartingY uint16

    PanSpeed uint16
}

type GraphicsConfig struct {
    PanningOptions

    ScreenX uint16
    ScreenY uint16

    SquareSize uint16

    TextureDir string
    FontFile string
    FontSize int
}
