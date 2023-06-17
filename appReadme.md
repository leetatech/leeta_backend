
<b>Run</b>
<ul>
    <li><i>cd cmd</i></li>
    <li><i>DATABASE_URL="postgres://postgres:9363T@localhost:5433/leeta?sslmode=disable" PORT="3000" go run main.go</i></li>
</ul>

<b>Swagger Doc</b>
<ul>
    <li>If you have initiated new interface implementation, you need to initialize the swagger doc again. Do this by running 
        <b><i>swag init --parseDependency --parseInternal -dir ../cmd -o ../docs</i></b>
    </li>

<li>
    <a href="http://localhost:3000/leeta/swagger/index.html">http://localhost:3000/leeta/swagger/index.html</a>
</li>
</ul>


<b>Application Health</b>
<p>Click <i><b><a href="http://localhost:3000/health">here (http://localhost:3000/health) </a></b></i> to test the running application.</p>

