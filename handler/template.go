package handler

import "html/template"

var defaultTemplate = template.Must(template.New("template").Parse(`
<html>
	<head>
	<script src="https://cdn.jsdelivr.net/npm/@json-editor/json-editor@latest/dist/jsoneditor.min.js"></script>
	</head>
	<body>

		<div id="editor_holder">
		</div>
		<button id="btn-submit">Submit</button>

		<script>
			(() => {
				var element = document.getElementById('editor_holder');
				var editor = new JSONEditor(element, {
					schema: JSON.parse("{{.Schema}}")
				});
				{{if .JSON}}editor.setValue(JSON.parse("{{.JSON}}"));{{end}}
				var btnSubmit = document.getElementById('btn-submit');
				btnSubmit.addEventListener('click', () => { 
					fetch(".", {
						method: "POST",
						headers: {
							"Content-Type": "application/json; charset=utf-8",
						},
						body: JSON.stringify(editor.getValue())
					}).catch(error => console.error(error));
				});
			})();
		</script>
	</body>
</html>`))
