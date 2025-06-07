package train

import (
	"fmt"
	stl2 "github.com/timleecasey/stllib/lib/stl"
	"math/rand"
	"time"

	torch "github.com/wangkuiyi/gotorch"
	"github.com/wangkuiyi/gotorch/nn"
	"github.com/wangkuiyi/gotorch/nn/functional"
)

// LoadSTL reads an STL file and extracts its mesh features into a tensor.
func LoadSTL(filename string) (*torch.Tensor, error) {
	//file, err := os.Open(filename)
	//if err != nil {
	//	log.Fatalf("Failed to open STL file: %v", err)
	//}
	//defer file.Close()
	//
	//model, err := stl.ReadFile(file)
	//if err != nil {
	//	log.Fatalf("Failed to read STL file: %v", err)
	//}
	//
	//data := make([]float32, numTriangles*9) // Each triangle has 3 vertices with 3 coordinates each
	//
	//for i, tri := range model.Triangles {
	//	for j := 0; j < 3; j++ {
	//		data[i*9+j*3+0] = tri.Vertices[j][0] // X
	//		data[i*9+j*3+1] = tri.Vertices[j][1] // Y
	//		data[i*9+j*3+2] = tri.Vertices[j][2] // Z
	//	}
	//}

	stl, err := stl2.LoadModel("example.stl")
	if err != nil {
		return nil, err
	}
	numTriangles := len(*stl.Objs)
	data := make([]float32, numTriangles*9) // Each triangle has 3 vertices with 3 coordinates each
	for i := 0; i < numTriangles; i++ {
		j := 0
		data[i*9+j*3+0] = float32((*stl.Objs)[i].A.X) // X
		data[i*9+j*3+1] = float32((*stl.Objs)[i].A.Y) // Y
		data[i*9+j*3+2] = float32((*stl.Objs)[i].A.Z) // Z

		j = 1
		data[i*9+j*3+0] = float32((*stl.Objs)[i].B.X) // X
		data[i*9+j*3+1] = float32((*stl.Objs)[i].B.Y) // Y
		data[i*9+j*3+2] = float32((*stl.Objs)[i].B.Z) // Z

		j = 2
		data[i*9+j*3+0] = float32((*stl.Objs)[i].C.X) // X
		data[i*9+j*3+1] = float32((*stl.Objs)[i].C.Y) // Y
		data[i*9+j*3+2] = float32((*stl.Objs)[i].C.Z) // Z

	}

	fmt.Printf("Extracted %d triangles from STL file.\n", numTriangles)
	return &torch.NewTensor(data).Reshape([]int64{1, int64(len(data))}), nil
}

// ToolpathModel defines a simple neural network for STL-to-toolpath mapping.
type ToolpathModel struct {
	net *nn.SequentialModule
}

// NewToolpathModel initializes the neural network model.
func NewToolpathModel() *ToolpathModel {
	model := &ToolpathModel{
		net: nn.Sequential(
			nn.Linear(9*256, 512, false),
			functional.Relu(),
			nn.Linear(512, 256, false),
			functional.Relu(),
			nn.Linear(256, 9*256, false),
		),
	}
	return model
}

func MSELoss(input, target torch.Tensor) torch.Tensor {
	diff := input.Sub(target)
	squaredDiff := diff.Mul(diff)
	return squaredDiff.Mean()
}

// TrainModel trains the neural network model using the provided STL data.
func TrainModel(model *ToolpathModel, stlData torch.Tensor) {
	rand.Seed(time.Now().UnixNano())
	optimizer := torch.SGD(model.net.Parameters(), 0.01, 0, 0, false)

	for epoch := 0; epoch < 100; epoch++ {
		// Simulate an ideal toolpath (denoised STL features)
		cleanToolpath := stlData.Mul(torch.NewTensor([]float32{0.95}))

		// Forward pass
		predictedToolpath := model.net.Forward(stlData)

		// Compute loss
		loss := MSELoss(predictedToolpath, cleanToolpath, torch.ReductionMean)
		optimizer.ZeroGrad()
		loss.Backward()
		optimizer.Step()

		fmt.Printf("Epoch %d - Loss: %v\n", epoch+1, loss.Item())
	}

	// Save the trained model
	torch.Save(model.net, "toolpath_model.pt")
	fmt.Println("Training complete! Model saved as toolpath_model.pt")
}

func main() {
	// Load STL data
	stlData := LoadSTL("example.stl")

	// Initialize and train the model
	model := NewToolpathModel()
	TrainModel(model, stlData)
}
