package train

import (
	"fmt"
	"github.com/timleecasey/stllib/lib/stl"
	torch "github.com/wangkuiyi/gotorch"
	"github.com/wangkuiyi/gotorch/nn"
	"github.com/wangkuiyi/gotorch/nn/functional"
	"log"
	"os"
)


// LoadSTL reads an STL file and extracts its mesh features into a tensor.
func LoadSTL(filename string) torch.Tensor {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open STL file: %v", err)
	}
	defer file.Close()

	model, err := stl.ReadFile(file)
	if err != nil {
		log.Fatalf("Failed to read STL file: %v", err)
	}

	numTriangles := len(model.Triangles)
	data := make([]float32, numTriangles*9) // Each triangle has 3 vertices with 3 coordinates each

	for i, tri := range model.Triangles {
		for j := 0; j < 3; j++ {
			data[i*9+j*3+0] = tri.Vertices[j][0] // X
			data[i*9+j*3+1] = tri.Vertices[j][1] // Y
			data[i*9+j*3+2] = tri.Vertices[j][2] // Z
		}
	}

	fmt.Printf("Extracted %d triangles from STL file.\n", numTriangles)
	return torch.NewTensor(data).Reshape([]int64{1, int64(len(data))})
}

// LoadModel loads the trained neural network model from a file.
func LoadModel(filepath string) *nn.SequentialModule {
	model := nn.Sequential(
		nn.Linear(9*256, 512, false),
		nn.ReLU(),
		nn.Linear(512, 256, false),
		nn.ReLU(),
		nn.Linear(256, 9*256, false),
	)
	err := torch.Load(model, filepath)
	if err != nil {
		log.Fatalf("Failed to load model: %v", err)
	}
	fmt.Println("Model loaded successfully.")
	return model
}

// GenerateToolpath uses the model to predict an optimized toolpath from STL data.
func GenerateToolpath(model *nn.SequentialModule, stlData torch.Tensor) torch.Tensor {
	return model.Forward(stlData)
}

// ConvertToGCode converts the toolpath tensor into G-code commands.
func ConvertToGCode(toolpath torch.Tensor) string {
	data := toolpath.ToData().([]float32)
	gcode := "G21 ; Set units to mm\nG90 ; Absolute positioning\n"
	for i := 0; i < len(data); i += 3 {
		gcode
		::contentReference[oaicite:0]{index=0}


