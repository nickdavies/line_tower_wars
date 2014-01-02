package texture

import (
    "errors"
    "strconv"
    "fmt"
    "strings"
    "os"
    "bufio"
    "path"
)

import (
    "github.com/neagix/Go-SDL/sdl"
)

const texture_list_filename = "textures.txt"

var NoSuchTextureErr = errors.New("No such texture exists")

type Texture struct {
    Id int
    Name string
    Height uint16
    Width uint16
    Surface *sdl.Surface
}

type TextureMap interface {
    LookupId(texture_name string) (texture_id int, err error)

    GetName(texture_name string) *Texture
    GetId(texture_id int) *Texture
}

func NewTextureMap(square_size uint16, texture_dir string) (TextureMap, error) {

    textures, texture_names, err := loadTextures(texture_dir)
    if err != nil {
        return nil, err
    }

    return &textureMapStruct{
        textures: textures,
        texture_names: texture_names,
    }, nil
}

type textureMapStruct struct {
   textures map[int]*Texture
   texture_names map[string]int
}

func (tm *textureMapStruct) LookupId(texture_name string) (int, error) {
    id, ok := tm.texture_names[texture_name]
    if ok {
        return id, nil
    }
    return 0, NoSuchTextureErr
}

func (tm *textureMapStruct) GetName(texture_name string) *Texture {
    id, ok := tm.texture_names[texture_name]
    if !ok {
        panic(NoSuchTextureErr)
    }

    texture, ok := tm.textures[id]
    if !ok {
        panic(fmt.Errorf("Texture with id %d in texture_names but not in textures", id))
    }

    return texture
}

func (tm *textureMapStruct) GetId(texture_id int) *Texture {
    texture, ok := tm.textures[texture_id]
    if !ok {
        panic(NoSuchTextureErr)
    }

    return texture
}

// Texture loading functions

func loadTextures(texture_dir string) (textures map[int]*Texture, texture_names map[string]int, err error) {

    f, err := os.Open(path.Join(texture_dir, texture_list_filename))
    if err != nil {
        return nil, nil, err
    }

    textures = make(map[int]*Texture, 0)
    texture_names = make(map[string]int)

    line_num := 0
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line_num++
        line := strings.TrimSpace(scanner.Text())
        if len(line) == 0 {
            continue
        }

        fields := strings.Split(line, ",")

        if len(fields) != 5 {
            return nil, nil, fmt.Errorf("Error on line %d: wrong number of fields got %d expected 5", line_num, len(fields))
        }

        filename := strings.TrimSpace(fields[4])
        name := strings.TrimSpace(fields[3])
        id, err := strconv.ParseInt(strings.TrimSpace(fields[0]), 10, 32)
        if err != nil {
            return nil, nil, fmt.Errorf("Error on line %d: texture id is invalid", line_num)
        }
        width, err := strconv.ParseInt(strings.TrimSpace(fields[1]), 10, 32)
        if err != nil {
            return nil, nil, fmt.Errorf("Error on line %d: texture width is invalid", line_num)
        }
        height, err := strconv.ParseInt(strings.TrimSpace(fields[2]), 10, 32)
        if err != nil {
            return nil, nil, fmt.Errorf("Error on line %d: texture height is invalid", line_num)
        }

        texture_surface, err := loadTexture(path.Join(texture_dir, filename))
        if err != nil {
            return nil, nil, fmt.Errorf("Error on line %d: failed to load file (%s): %s", line_num, filename, err)
        }

        _, ok := textures[int(id)]
        if ok {
            return nil, nil, fmt.Errorf("Error on line %d: duplicate id %d", line_num, id)
        }
        _, ok = texture_names[name]
        if ok {
            return nil, nil, fmt.Errorf("Error on line %d: duplicate name %s", line_num, name)
        }

        textures[int(id)] = &Texture{
            Id: int(id),
            Name: name,
            Width: uint16(width),
            Height: uint16(height),
            Surface: texture_surface,
        }
        texture_names[name] = int(id)
    }

    return textures, texture_names, nil
}

func loadTexture(texture_file string) (*sdl.Surface, error) {

    surface := sdl.Load(texture_file)
    if surface == nil {
        return nil, fmt.Errorf("Failed to load texture: %s", sdl.GetError())
    }

    surface = surface.DisplayFormat()
    if surface == nil {
        return nil, fmt.Errorf("Failed to convert texture: %s", sdl.GetError())
    }

    return surface, nil
}

