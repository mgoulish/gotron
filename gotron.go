package main

import ( 
         "fmt"
         "io"
         "net"
         "os"
         "time"
       )



var fp = fmt.Fprintf


var total_bytes int



func cnx_handler ( name       string,
                   cnx_number int, 
                   cnx        net.Conn ) {

  buffer := make ( []byte, 512 )

  for {
    var   n   int
    var err error

    n, err = cnx.Read ( buffer )

    if err != nil {
      if err != io.EOF {
        fp ( os.Stderr, "%s error : connection read : |%s|\n", name, err.Error() )
      }
      break
    }

    if n > 0 {
      message := string(buffer[0:n])
      fp ( os.Stdout, "%s : received from cnx %d : |%s|\n", name, cnx_number, message )
      total_bytes += len ( message )
      fp ( os.Stdout, "%s has received %d total bytes.\n", name, total_bytes )
    }
    
    //experiment 
    cnx.Write ( buffer )
  }
}





// This can accept multiple connections, and will 
// launch a goroutine for each.
func listen ( name string, port string ) {
  cnx_count := 0

  tcp_listener, err := net.Listen ( "tcp", ":" + port )
  if err != nil {
    fp ( os.Stdout, "name error : net.Listen error |%s|\n", name, err.Error() )
    os.Exit ( 1 )
  }

  for {
    cnx, err := tcp_listener.Accept ( )

    if err != nil {
      fp ( os.Stdout, "name error : |%s|\n", name, err.Error() )
      os.Exit(1)
    }
    fp ( os.Stdout, "accept returned.\n" )

    cnx_count ++
    go cnx_handler ( name, cnx_count, cnx )
  }
}





func dialer ( name string, host, port string, timeout_seconds int ) ( cnx net.Conn, err error ) {

  for t := 0; t < timeout_seconds; t ++ {
    cnx, err := net.Dial ( "tcp", host + ":" + port )
    if err == nil {
      return cnx, nil
    }

    time.Sleep ( time.Second )
  }

  return nil, fmt.Errorf ( "%s : timed out", name )
}





func make_connection ( name string, n_messages int, host, port string ) {
  cnx, err := dialer ( name, host, port, 10 )

  if err != nil {
    fp ( os.Stdout, "%s error: |%s|\n", name, err.Error() )
    os.Exit ( 1 )
  }

  if cnx == nil {
    fp ( os.Stdout, "%s error : nil connection.\n", name )
    os.Exit ( 2 )
  }

  defer cnx.Close ( )

  message := "0123456789"

  for i := 0; i < n_messages; i ++ {
    fp ( cnx, message )
    fp ( os.Stdout, "%s sent |%s|\n", name, message )
    time.Sleep ( 5 * time.Second )
  }
}





func main ( ) {
  name          := os.Args[1]
  host          := os.Args[2]
  incoming_port := os.Args[3]
  outgoing_port := os.Args[4]

  go listen ( name, incoming_port )

  n_messages := 10
  go make_connection ( name, n_messages, host, outgoing_port )


  for {
    time.Sleep ( 100 * time.Second )
  }
}





