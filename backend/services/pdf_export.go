package services

import (
	"bytes"
	"fmt"

	"github.com/fastlearn/backend/assets"
	"github.com/fastlearn/backend/models"
	"github.com/go-pdf/fpdf"
)

// generateSummaryPDF renders a video's summary (overview, key points, takeaways,
// examples, Q&A) into a PDF. DejaVu Sans is embedded so Cyrillic/Kazakh text renders correctly.
func generateSummaryPDF(video *models.Video, summary *models.Summary) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8FontFromBytes("DejaVu", "", assets.DejaVuSansRegular)
	pdf.AddUTF8FontFromBytes("DejaVu", "B", assets.DejaVuSansBold)
	pdf.SetMargins(18, 18, 18)
	pdf.AddPage()

	title := "Video Summary"
	if video.Title != nil && *video.Title != "" {
		title = *video.Title
	}
	pdf.SetFont("DejaVu", "B", 18)
	pdf.MultiCell(0, 8, title, "", "L", false)

	pdf.SetFont("DejaVu", "", 10)
	pdf.SetTextColor(120, 120, 120)
	meta := summary.Language
	if video.DurationSeconds != nil {
		meta = fmt.Sprintf("%s · %d min", meta, (*video.DurationSeconds+59)/60)
	}
	pdf.MultiCell(0, 6, meta, "", "L", false)
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(4)

	if summary.Content != nil {
		c := summary.Content
		if c.Overview != "" {
			pdf.SetFont("DejaVu", "", 11)
			pdf.MultiCell(0, 6, c.Overview, "", "L", false)
			pdf.Ln(4)
		}

		if len(c.KeyPoints) > 0 {
			writeSectionHeading(pdf, "Key Points")
			for _, kp := range c.KeyPoints {
				pdf.SetFont("DejaVu", "B", 11)
				pdf.MultiCell(0, 6, kp.Topic, "", "L", false)
				pdf.SetFont("DejaVu", "", 10)
				pdf.MultiCell(0, 6, kp.Explanation, "", "L", false)
				pdf.Ln(2)
			}
			pdf.Ln(2)
		}

		if len(c.KeyTakeaways) > 0 {
			writeSectionHeading(pdf, "Key Takeaways")
			pdf.SetFont("DejaVu", "", 10)
			for _, t := range c.KeyTakeaways {
				pdf.MultiCell(0, 6, "• "+t, "", "L", false)
			}
			pdf.Ln(2)
		}

		if len(c.Examples) > 0 {
			writeSectionHeading(pdf, "Examples")
			for _, ex := range c.Examples {
				pdf.SetFont("DejaVu", "B", 11)
				pdf.MultiCell(0, 6, ex.Topic, "", "L", false)
				pdf.SetFont("DejaVu", "", 10)
				pdf.MultiCell(0, 6, ex.Example, "", "L", false)
				pdf.Ln(2)
			}
			pdf.Ln(2)
		}
	}

	if len(summary.QA) > 0 {
		writeSectionHeading(pdf, "Q&A")
		for _, qa := range summary.QA {
			pdf.SetFont("DejaVu", "B", 10)
			pdf.MultiCell(0, 6, "Q: "+qa.Question, "", "L", false)
			pdf.SetFont("DejaVu", "", 10)
			pdf.MultiCell(0, 6, "A: "+qa.Answer, "", "L", false)
			pdf.Ln(2)
		}
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func writeSectionHeading(pdf *fpdf.Fpdf, text string) {
	pdf.SetFont("DejaVu", "B", 13)
	pdf.MultiCell(0, 7, text, "", "L", false)
	pdf.Ln(1)
}
