package config

func fiasConfig(m map[string]string) {
	m["ARCHIVE_LOCAL_PATH"] = ""
	m["IS_NEED_DOWNLOAD_ARCHIVE"] = ""
	m["ARCHIVE_PAGE_LINK"] = "https://fias.nalog.ru/Frontend"
	m["ARCHIVE_LINK_SELECTOR"] = "#OpenDataPagerDataPlaceholder a"
	m["OBJECT_FILE_PART"] = "AS_ADDR_OBJ"
	m["HOUSES_FILE_PART"] = "AS_HOUSES_"
	m["HIERARCHY_FILE_PART"] = "_HIERARCHY_"
	m["IMPORT_DESTINATION"] = ""
}
