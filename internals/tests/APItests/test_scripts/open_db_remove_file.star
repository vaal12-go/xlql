# This starlark file will create and populate (this is needed)
# a db file, which then will be tested to be removed by next run 
# of this script. This will test the 'delete_db_if_exists' parameter 
# of open_db function.

hrDB = open_db(
    "hr_db.sqlite"
    ,delete_db_if_exists = True
)

countries_q = hrDB.load_excel_sheet(
    file_name = "../../../examples/Sample HR Database.xlsx",
    sheet_name = "countries"
)
print("DB file should have been created")



