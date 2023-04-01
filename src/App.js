import './App.css';
import { useEffect, useRef, useState } from 'react';

const AjaxFetch = (url, obj, method) => fetch(url, {
  method: method,
  credentionals: 'include',
  headres: {
    'Accept': 'application/json',
    'Conntent-Type': 'application/json'
  },
  body: JSON.stringify(obj)
}).then(function (response) {
  if (response.status !== 200) {
    console.log(`fetch return not ok ${response.status}`);
  }
  return response.json().then((result) => {
    return result
  })
})
const Tooltip = ({ children, text, ...rest }) => {
  // const [show, setShow] = useState(false);

  return (
    <div className="tooltip">
      <div
        // onMouseEnter={() => setShow(true)}
        // onMouseLeave={() => setShow(false)}
        {...rest}
      >
        {children}
      </div>
      {(text != "") ? <div className="tooltip-box">
        {text}
        {/* <span className="tooltip-arrow" /> */}
      </div> : <></>}
    </div>
  );
};
function Parse(data) {
  const ErrReg = new RegExp(`(ERROR)`)
  const WarnReg = new RegExp(`(WARN)`)
  var err_data = data.split(ErrReg)
  var warn_data = err_data.map(err => err.split(WarnReg)).flat()
  return warn_data
}
function App() {
  const ref = useRef()
  const [getdata, setdata] = useState([])
  const [getselect, setselect] = useState("")
  const [dialog, setdialog] = useState("hidden")
  const [loggs, setloggs] = useState([])
  var col
  //   ip: "",
  //   list_port: [{ port: "", status: "" }],
  //   list_url: [{ url: "", status: "" }],
  //   list_ssh: [{ name: "", result: "" }],
  //   status:[{name: "", result: ""}]
  // }])
  useEffect(() => {
    AjaxFetch('/api/get', {}, 'POST').then((data) => {
      console.log(data)
      setdata(data)
      AjaxFetch('/api/errorlogs', {}, 'POST').then((logerr) => {
        console.log(logerr)
        let err = data.Stend.map(host => {
          console.log(host)
          let log = logerr.filter(logs => host.ip == logs.host)
          console.log(log)
          if (log.length > 0) {
            host.status.map(st => {
              err = log.find(l => l.module == st.status)
              console.log("Err", err)
              if (err) {
                st.errors = err.errors
              }
              return st
            })
          }
          return host
        })
        console.log(err)
        setdata({list_url: data.list_url, Stend: err })
      })
    })
  }, []);
  useEffect(() => {
    const checkIfClickedOutside = e => {
      if (dialog == "visible" && ref.current && !ref.current.contains(e.target)) {
        setdialog('hidden')
      }
    }
    document.addEventListener("mousedown", checkIfClickedOutside)
    return () => {
      document.removeEventListener("mousedown", checkIfClickedOutside)
    }
  }, [dialog])
  const onHover = ip => {
    setselect(ip)
  }
  const sshconnect = ip => {
    AjaxFetch('/api/terminal', { ip: ip }, 'POST').then(() => { })
  }
  const Warlogs = (host, module, service) => {
    AjaxFetch('/api/warlog', { host: host, service: service, module: module }, 'POST').then((data) => {
      console.log(data)
      if (!data) return
      setloggs(Parse(data))
      setdialog("visible")
    })
  }
  const DialogBox = () => {
    setdialog("visible")
  }
  return (
    (getdata.length == 0) ? <div><h1>Loading..</h1></div> : <>
      <div className='modal' style={{ visibility: dialog }} >
        <div ref={ref} className='modal-content'><pre>
          {loggs.map(log => (log == "ERROR") ? <span style={{backgroundColor: "red"}}>{log}</span> : 
          (log == "WARN") ? <span style={{backgroundColor: "orange"}}>{log}</span> :<span>{log}</span> )}
          </pre>
          {/* <pre>{loggs}</pre> */}
        </div>
      </div>
      <div className='parent'>
        <div className='div1'>
          <table>
            {getdata.Stend.map(host => {
              // (host.ip == getselect) ? col = "red" : col = "black"
              return <tr onMouseEnter={() => onHover(host.ip)} className={(host.ip == getselect) ? 'sele' : ""}>
                <td onClick={() => sshconnect(host.ip)}>
                  <Tooltip text={host.ip}>
                    {host.name}
                  </Tooltip>
                </td>
                {host.list_port.map(port => {
                  var color = (port.status == "failed") ? "red" : "rgb(130, 244, 160)"
                  return <>
                    <td style={{ backgroundColor: color }}>{port.port}</td>
                  </>
                })}
              </tr>
            })}
          </table>
        </div>
        <div className='div2'>
          <table>
            {getdata.Stend.map(host => {
              // (host.ip == getselect) ? col = "red" : col = "black"
              return <tr onMouseEnter={() => onHover(host.ip)} className={(host.ip == getselect) ? 'sele' : ""}>
                {(host.list_ssh.length == 0) ? <td>...</td> : host.list_ssh.map(s => {
                  var color = (s.result == "active\n" || s.result.includes("is running")) ? "rgb(130, 244, 160)" : "red"
                  color = (s.result == "Timeout") ? "red" : color
                  if (s.name == "DiskFree") {
                    return <td>{s.result}</td>
                  }
                  if (s.name == "LoadAvg") {
                    return <td>{s.result}</td>
                  }
                  return <>
                    <td className='tdcuttext' style={{ backgroundColor: color }}>{s.name}</td>
                  </>
                })}
              </tr>
            })}
          </table>
        </div>
        <div className='div4'>
          <table>
            <tr>
              {getdata.list_url.map(url => {
                // (host.ip == getselect) ? col = "red" : col = "black"
                var color = (url.status == "200") ? "rgb(130, 244, 160)" : "red"
                return <>
                  <td style={{ backgroundColor: color }}>{url.url}</td>
                </>
              })}
            </tr>
          </table >
          {getdata.Stend.map(host => (host.status) ? <table>
            <tr onMouseEnter={() => onHover(host.ip)} className={(host.ip == getselect) ? 'sele' : ""}>
              {/* <tr onMouseEnter={() => onHover(host.ip)}> */}
              <td onClick={() => sshconnect(host.ip)}>
                <Tooltip text={host.ip}>
                  {host.name}
                </Tooltip>
              </td>
              {host.status.map(stat => {
                var color = "rgb(130, 244, 160)"
                var color = (stat.errors && stat.errors > 0) ? "orange": color
                color = (stat.result == "running") ? color : "red"
                color = (stat.result == "Timeout") ? "red" : color
                return <td style={{ backgroundColor: color }}
                  onClick={() => Warlogs(host.ip, stat.status, stat.service)}>
                  <Tooltip text={stat.tooltip}>
                    <b>{(stat.errors && stat.errors != 0) ? stat.errors : ""}</b>  {stat.status}
                  </Tooltip>
                </td>
              })}
            </tr>
          </table> : <></>)}
        </div>
      </div >
    </>
  );
}

export default App;
