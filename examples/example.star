hrDB = open_db(
    "hr_db.sqlite"
    ,delete_db_if_exists = True
) #[x]: check if DB can be deleted upon opening (if exists)

countries_q = hrDB.load_excel_sheet(
    file_name = "Sample HR Database.xlsx",
    sheet_name = "countries"
)

countries_q.print()