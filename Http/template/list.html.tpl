<!DOCTYPE html>
<html>
<head>
   <title>test 1</title>
   <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
   <link href="http://libs.baidu.com/bootstrap/3.0.3/css/bootstrap.min.css" rel="stylesheet">
   <script src="http://libs.baidu.com/jquery/2.0.0/jquery.min.js"></script>
   <script src="http://libs.baidu.com/bootstrap/3.0.3/js/bootstrap.min.js"></script>
<style type="text/css">
</style>
</head>
<body>
          <div class="col-md-6">
            <table class="table table-striped table-bordered dataTable no-footer">
              <thead>
                <th>ID</th>
                <th>名字</th>
                <th>备注</th>
							  <th>是否admin</th>
							  <th>操作</th>
							  <th>操作</th>
              </thead>
              <tbody id='table-content'>
						    {{range .}}
						    <tr>
							    <td>{{.Id}}</td>
							    <td>{{.Name}}</td>
							    <td>{{.Note}}</td>
							    <td>{{.Isadmin}}</td>
							    <td><a href='/delete?id={{.Id}}'>删除</a> </td>
                  <td><a href='/update?id={{.Id}}'>更新</a> </td>
							  </tr>
						    {{end}}
              </tbody>
            </table>
          </div>
		</br>
		<div class="col-md-10">
		  <button class="btn btn-primary"> 新增</button>
		</div>

<script type="text/javascript">
</script>
</body>
</html>