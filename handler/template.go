package handler

import "html/template"

// Template is default template which creates an editor for a JSON Schema.
var Template = template.Must(template.New("template").Parse(`
<html>
	<head>
		<title>JSON Editor</title>
		<script src="https://cdn.jsdelivr.net/npm/@json-editor/json-editor@latest/dist/jsoneditor.min.js"></script>
	</head>
	<body>

		<div id="editor_holder">
		</div>
		<button id="btn-submit">Submit</button>

		<script>
			var element = document.getElementById('editor_holder');
			var editor = new JSONEditor(element, {
				schema: JSON.parse("{{.Schema}}")
			});
			{{if .JSON}}editor.setValue(JSON.parse("{{.JSON}}"));{{end}}

			function defaultSubmit(val) {
				fetch(".", {
					method: "POST",
					headers: {
						"Content-Type": "application/json; charset=utf-8",
					},
					body: val 
				}).catch(error => console.error(error));
			}

			var btnSubmit = document.getElementById('btn-submit');
			btnSubmit.addEventListener('click', () => {
				var val = JSON.stringify(editor.getValue());
				if (typeof submit === "undefined") {
					defaultsubmit(val);
				} else {
					submit(val);
				}
			});
		</script>
	</body>
</html>`))
