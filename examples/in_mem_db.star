
#In memory example database
#More info: https://www.sqlite.org/inmemorydb.html
inMemDB = open_db("file:memdb1?mode=memory&cache=shared")

mem_tbl = inMemDB.create_table(
    name = "qwe1",
    columns = {
	     "Col1": "TEXT",
	     "Col23": "NUMERIC",
	     "Col2": "NUMERIC",
	     "Col0": "NUMERIC"
	}
)

mem_tbl.add_row("qwe1", 1, 2, 3)

mem_q = inMemDB.run_query("SELECT * FROM qwe1")

mem_q.print()