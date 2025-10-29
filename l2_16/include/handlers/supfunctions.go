package handlers

import "golang.org/x/net/html"

func extractAttr(attrs []html.Attribute, key string) string {
	for _, attr := range attrs {
		if attr.Key == key {
			return attr.Val
		}
	}

	return ""
}
 
// looking for certain tag
func findCertainTagAndItsData(n *html.Node, linksStorages *htmlIncludedLinks) {
	if n.Type == html.ElementNode { // is tag and isnt data included opening and closing tags
		if n.Data == "link" {
			rel := extractAttr(n.Attr, "rel")
			href := extractAttr(n.Attr, "href")

			// isHrefExist
			if rel == "stylesheet" && href != "" {
				linksStorages.CssLinks = append(linksStorages.CssLinks, href)
			}
		} else if n.Data == "script" {
			src := extractAttr(n.Attr, "src")
			if src != "" {
				linksStorages.JsLinks = append(*&linksStorages.JsLinks, src)
			}
		} else if n.Data == "img" {
			src := extractAttr(n.Attr, "srs")
			if src != "" {
				linksStorages.JsLinks = append(*&linksStorages.JsLinks, src)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findCertainTagAndItsData(c, linksStorages)
	}
}