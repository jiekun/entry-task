import time
import pymysql.cursors

conn = pymysql.connect(
   host='127.0.0.1',
   user='root',
   password='',
   database='entry_task_v2_db',
   charset='utf8mb4',
   cursorclass=pymysql.cursors.DictCursor
)

cursor = conn.cursor()

sql = """INSERT INTO user_tab (name, password, nickname, profile_pic, status, ctime, mtime) 
VALUES
("stress_test_{}", "49069a895b5743c929c8578b736fc015", "stress_test_nickname", "https://images2.imgbox.com/40/1e/n2bhfC9o_o.jpeg", 0, {}, {});
"""

print(sql)
for i in range(10000000, 20000000):
   try:
      tmp_sql = sql.format(i, int(time.time()), int(time.time()))
      cursor.execute(tmp_sql)
      if not i % 10000:
         print(i)
         print(tmp_sql)
         conn.commit()

   except Exception as e:
      print(e)

# Closing the connection
conn.commit()
conn.close()
