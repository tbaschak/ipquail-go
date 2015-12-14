<!DOCTYPE html>
<!--[if lt IE 7]>      <html class="no-js lt-ie9 lt-ie8 lt-ie7"> <![endif]-->
<!--[if IE 7]>         <html class="no-js lt-ie9 lt-ie8"> <![endif]-->
<!--[if IE 8]>         <html class="no-js lt-ie9"> <![endif]-->
<!--[if gt IE 8]><!--> <html class="no-js"> <!--<![endif]-->
    <head>
        <meta charset="utf-8">
        <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
        <title>IP Quail - What is My IP Address?</title>
        <meta name="description" content="">
        <meta name="viewport" content="width=device-width, initial-scale=1">

        <link href="css/bootstrap.min.css" rel="stylesheet">
        <style>
            body {
                padding-top: 50px;
                padding-bottom: 20px;
            }
        </style>

        <script src="js/vendor/modernizr-2.6.2-respond-1.1.0.min.js"></script>
    </head>
    <body>
        <!--[if lt IE 7]>
            <p class="browsehappy">You are using an <strong>outdated</strong> browser. Please <a href="http://browsehappy.com/">upgrade your browser</a> to improve your experience.</p>
        <![endif]-->

    <!-- Main jumbotron for a primary marketing message or call to action -->
    <div class="jumbotron">
      <div class="container">
        <div class="row">
            <div class="col-xs-8 col-sm-8">
                <h1>IP Quail!</h1>
                <p>This website will tell you your IP[v4|v6] addresses.</p>
                <h3>Your IPv6 address is:</h3>
                <p id="ipv6">None</p>
                <h3>Your IPv4 address is:</h3>
                <p id="ipv4">None</p>
                <h3>Your IPv6 PTR is:</h3>
                <p id="ipv6ptr">None</p>
                <h3>Your IPv4 PTR is:</h3>
                <p id="ipv4ptr">None</p>
                <hr>
                <p>This website can also be shell scripted using curl:<br><code>IP6=`curl -s 6.henchcdn.com/ip`</code> or <code>IP4=`curl -s 4.henchcdn.com/ip`</code> or <code>PTR6=`curl -s 6.henchcdn.com/ptr`</code> or <code>PTR4=`curl -s 4.henchcdn.com/ptr`</code></p>
            </div>
            <div class="col-xs-4 col-sm-4">
                <img class="img-responsive" src="assets/img/bird.png" alt="A Funny Bird"/>
            </div>
        </div>
      </div>
    </div>

    <div class="container">
      <footer>
        <p><a href="https://github.com/tbaschak/">&copy; Theo Baschak 2015</a></p>
      </footer>
    </div> <!-- /container -->
    <script src="//ajax.googleapis.com/ajax/libs/jquery/1.11.0/jquery.min.js"></script>
    <script>window.jQuery || document.write('<script src="js/vendor/jquery-1.11.0.min.js"><\/script>')</script>

    <script src="js/vendor/bootstrap.min.js"></script>

    <script type="text/javascript" language="javascript">
    $(document).ready(function() {
      $.getJSON('http://6.henchcdn.com/api/ip', function(ab) {
        $('#ipv6').html('<code>' + ab.ip + '</code>');
      });
      $.getJSON('http://4.henchcdn.com/api/ip', function(ac) {
        $('#ipv4').html('<code>' + ac.ip + '</code>');
      });
      $.getJSON('http://6.henchcdn.com/api/ptr', function(ad) {
        $('#ipv6ptr').html('<code>' + ad.ptr + '</code>');
      });
      $.getJSON('http://4.henchcdn.com/api/ptr', function(ae) {
        $('#ipv4ptr').html('<code>' + ae.ptr + '</code>');
      });
    });
    </script>
    
    </body>
</html>
