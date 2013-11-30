using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;
using System.Text;
using System.Threading.Tasks;

namespace famousfront.utils
{
  [AttributeUsage(AttributeTargets.Property, AllowMultiple=false)]
  class HttpQueryParameter : Attribute
  {
    internal enum ParamType
    {
      Path,Query,
    }
    readonly string _name;
    readonly ParamType _type;
    public string Name { get { return _name; } }
    public ParamType Type { get { return _type; } }

    internal HttpQueryParameter()
      : this(null, ParamType.Query)
    {

    }
    internal HttpQueryParameter(string name)
      : this(name, ParamType.Query)
    {

    }
    internal HttpQueryParameter(string name, ParamType t)
    {
      _name = name;
      _type = t;
    }
    internal static void EnumFields(object request, Action<string, object> action)
    {
      foreach (var field in request.GetType().GetProperties(BindingFlags.Instance | BindingFlags.Public))
      {
        var name = field.Name.ToLower();
        var val = field.GetValue(request);
        if (val != null)
        {
          action(name, val);
        }
      }
    }
    internal static string EncodeFieldQueryString(object request)
    {
      var queryparams = new Dictionary<string, string>();
      EnumFields(request, (n, v) =>
      {
        var strv = v.ToString();
        if (!string.IsNullOrEmpty(strv))
        {
          queryparams.Add(n, strv);
        }
      });
      var vals = new List<string>();
      foreach (var p in queryparams)
      {
        vals.Add(String.Format("{0}={1}", p.Key, Uri.EscapeDataString(p.Value)));
      }
      return string.Join("&", vals);
    }
    internal static void EnumParams(object request, Action<ParamType, string, object> action)
    {
      foreach (var prop in request.GetType().GetProperties(BindingFlags.Instance | BindingFlags.Public))
      {
        var attr = prop.GetCustomAttribute<HttpQueryParameter>();
        if (attr == null)
          continue;
        var name = attr.Name ?? prop.Name.ToLower();
        var pt = attr.Type;
        var val = prop.GetValue(request);
        if (prop.PropertyType.IsValueType || val != null)
        {
          action(pt, name, val);
        }
      }
    }
    internal static string EncodeQueryString(object request)
    {
      var queryparams = new List<KeyValuePair<string, string>>();
      EnumParams(request, (t, n, v) =>
      {
        var strv = v.ToString();
        if (t == ParamType.Query && !string.IsNullOrEmpty(strv))
        {
          queryparams.Add(new KeyValuePair<string, string>(n, strv));
        }
      });
      var vals = new List<string>();
      foreach (var p in queryparams)
      {
        vals.Add(String.Format("{0}={1}", p.Key, Uri.EscapeDataString(p.Value)));
      }
      return string.Join("&", vals);
    }
    internal static string EncodeRelativePath(object request, string pattern)
    {
      var pathparams = new Dictionary<string, string>();
      EnumParams(request, (t, n, v) => 
      {
        System.Diagnostics.Debug.Assert(string.IsNullOrEmpty(v.ToString()));
        if (t == ParamType.Path)
          pathparams.Add(n, v.ToString());
      });
      foreach (var p in pathparams)
      {
        var s = "{" + p.Key + "}";
        var v = Uri.EscapeDataString(p.Value);
        pattern = pattern.Replace(s, v);
      }
      return pattern;
    }
  }
  internal static class QueryStringEncoder
  {
    internal static string QueryString(this object request)
    {
      return HttpQueryParameter.EncodeFieldQueryString(request);
    }
  }
}
