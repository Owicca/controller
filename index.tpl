<html>
  <body>
    <ul id="list">
      {{range $key, $file := .}}
        <li>
          <a href="/v/{{$key}}/" target="_blank">
            <b>{{$file.Name}}</b>
          </a>
          <button class="del">Delete</button>
          </li>
      {{end}}
    </ul>
    <script type="text/javascript">
      let list = document.getElementById("list");
      let buttons = document.querySelectorAll(".del");
      buttons.forEach((but) => {
          but.addEventListener("click", (e) => {
            let url = e.target.previousElementSibling.getAttribute("href");
            fetch(url, {
              method: 'POST',
              headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json',
              },
              body: JSON.stringify({
                firstParam: 'Value',
                secondParam: 'OtherValue',
              })
            }).then((res) => {
              let data = res.json();
              if(data.success == true) {
                  e.target.parentNode.remove();
                } else {
                  alert(data.error);
                }
            });
          });
      });
    </script>
  </body>
</html>
