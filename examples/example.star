print("This is an example file")

hrDB = open_db("hr_db.sqlite")

countries_q = hrDB.load_excel_sheet(
    file_name = "Sample HR Database.xlsx",
    sheet_name = "countries"
)

countries_q.print()