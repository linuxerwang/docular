deb-spec {
	control {
		pkg-name: "docular"
		maintainer: "linuxerwang@gmail.com"
		description: "Docular doc reader server."

		other-attrs: {
			"Section": "utils",
			"Priority": "optional",
		}
	}

	content {
		path: "bin/docular"
		deb-path: "usr/bin/docular"
	}

	content {
		path: "webstatic"
		deb-path: "usr/share/docular/webstatic"
	}
}
