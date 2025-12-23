package formatter

import (
	"fmt"
	"io"
	"text/tabwriter"

	"gophkeep/internal/core/models"
)

type Formatter interface {
	FormatOutput(output io.Writer, secrets []models.SecretInfo)
}

type formatter struct {
}

func NewFormatter() Formatter {
	return &formatter{}
}

func (f *formatter) FormatOutput(output io.Writer, secrets []models.SecretInfo) {
	w := tabwriter.NewWriter(output, 10, 1, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "ID\tTYPE\tCREATED_AT\t")

	for _, secret := range secrets {
		fmt.Fprintf(w, "%d\t%s\t%s\t\n", secret.ID, secret.Type, secret.CreatedAt)
	}
}
