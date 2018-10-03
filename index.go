package main

var indexFile = []byte(`
<html>
<head>
    <title>Open garage door</title>
    <style>
        body{
            margin: 0;
            padding: 0;
            visibility: hidden;
            -webkit-touch-callout: none;
            -webkit-user-select: none;
            -khtml-user-select: none;
            -moz-user-select: none;
            -ms-user-select: none;
            -o-user-select: none;
            user-select: none;
        }
        .main_button {
            border-radius: 5px;
            border: 0;
            background-color: #3D9970;
            font-size: 17vmin;
            width: 100vw;
            height: 100vh;
            margin: 0;
        }
        .main_button.active {
            background-color: #FF4136;
            color: white;
        }
        .token {
            width: 100%;
            margin-top: 50px;
        }
    </style>
    <script>
    /* IE9+ ajax: atto({url: '/', method: 'GET', data: '', headers: {key: 'value'}}) .done(...) .error(...).finally(...) */
    !function(a){a.atto=function(a){var b,c,d,e=new XMLHttpRequest;e.open(a.method||"GET",a.url,!0),e.onload=function(){e.status>=200&&e.status<400?b&&b(e.responseText,e):c&&c(new Error(e.statusText),e),d&&d(null,e)};for(var f in a.headers)a.headers.hasOwnProperty(f)&&e.setRequestHeader(f,a.headers[f]);return e.onerror=function(){var a=new Error("Connection error");c&&c(a,{}),d&&d(a)},setTimeout(function(){e.send(a.data)}),{done:function(a){return b=a,this},error:function(a){return c=a,this},finally:function(a){return d=a,this}}}}(window);

    </script>
    </head>
<body>
    <button type="button" class="main_button">OPEN DOOR</button><br />
    <input type="text" class="token" placeholder="token">
    <script>
        var body = document.querySelector('body');
        var mainBtn = document.querySelector('.main_button');
        var tokenInput = document.querySelector('.token');
        var token = tokenInput.value = localStorage.getItem('token');
        var setText = function(text) {mainBtn.innerHTML = text;};
        mainBtn.addEventListener('click', function() {
            if (!token) {
                setText('Missing token!<br/>(scroll down)');
                return;
            }

            mainBtn.classList.add('active');
            setText('...sending...');

            atto({
                url: '/token',
                method: 'POST',
                data: '',
                headers: { token: token }
            })
            .done(function(r) {
                console.log(r)
                setText('Opening...');
                setTimeout(function() {
                    mainBtn.classList.remove('active');
                    setText('OPEN DOOR');
                },2250);
            })
            .error(function(error, response) {
                mainBtn.classList.remove('active');
                if (response.status == 401) {
                    setText('Wrong token!!!');
                } else {
                    alert('An error occured: ' + error.message);
                    setText('OPEN DOOR');
                }
            })
        });

        tokenInput.addEventListener('input', function(event) {
            localStorage.setItem('token', event.target.value);
            token = event.target.value;
        });

        body.style.visibility = 'visible';
    </script>
</body>
</html>
`)
