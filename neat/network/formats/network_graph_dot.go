package formats

import (
	"deepneat/neat/network"
	"io"

	"gonum.org/v1/gonum/graph/encoding/dot"
)

// WriteDOT is to write provided network graph using the GraphViz DOT encoding.
// See DOT Guide: https://www.graphviz.org/pdf/dotguide.pdf
func WriteDOT(w io.Writer, n *network.Network) error {
	data, err := dot.Marshal(n, n.Name, "", "")
	if err != nil {
		return err
	}
	if _, err = w.Write(data); err != nil {
		return err
	}
	return nil
}
