#Example of query.get_cursor iteration
#Based on in_mem_db.star

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
for x in range(5):
    mem_tbl.add_row("qwe1", x, x*10, x+1)

mem_q = inMemDB.run_query("SELECT * FROM qwe1")
print("query.Print representation")
mem_q.print()

print("Iterator example")
rowNo = 0
for row in mem_q.get_cursor():
    rowNo+=1
    print("Have row#:", rowNo)
    print(row)