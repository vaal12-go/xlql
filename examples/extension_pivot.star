#This requires loading of pivot extension:
#   https://github.com/jakethaw/pivot_vtab
#   Downloadable from https://sqlpkg.org/all/
# .dll of extension also available in sqliteextensions folder of this repository

dbName = "pivotDB.sqlite3"
newDB = open_db(dbName)


newDB.load_excel_sheet(file_name = "pivot_source_data.xlsx", 
    sheet_name = "r", drop_table=True)
newDB.load_excel_sheet(file_name = "pivot_source_data.xlsx",
    sheet_name = "x", drop_table=True)
newDB.load_excel_sheet(file_name = "pivot_source_data.xlsx", 
    sheet_name = "c", drop_table=True)

del_q = newDB.run_query(
    "DROP TABLE IF EXISTS pivot"
)

new_q = newDB.run_query(
    """
    CREATE VIRTUAL TABLE IF NOT EXISTS pivot USING pivot_vtab(
        (SELECT id r_id FROM r),        -- Pivot table key query
        (SELECT id c_id, name FROM c),  -- Pivot table column definition query
        (SELECT val                     -- Pivot query
            FROM x 
            WHERE r_id = ?1
            AND c_id = ?2)
    );
    """
)  #Because this is an table create statement, this cannot be iterated as query.

#For some reason direct SELECT after creation of pivot table locks the database.
# So things commented below should not be executed right here. They may be executed after SELECT from pivot3
# q = newDB.run_query(
#     "SELECT * FROM pivot"
# )
# q.print()

#New SELECT statement should be run to iterate through records
#Creating new table from the pivot view (which is created by extension) and this does not lock the database for reads
del_q = newDB.run_query(
    "DROP TABLE IF EXISTS pivot3"
)

new_q2 = newDB.run_query(
    """
    CREATE TABLE pivot3 AS 
	SELECT * FROM pivot
    """
)

new_q3 = newDB.run_query(
    "SELECT * FROM pivot3;"
)

new_q3.print()



