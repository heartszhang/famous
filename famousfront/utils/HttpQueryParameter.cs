using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;

namespace famousfront.utils
{
  [AttributeUsage(AttributeTargets.Property, AllowMultiple=false)]
  class HttpQueryParameter : Attribute
  {
    readonly string _name;
    public string Name { get { return _name; } }

    internal HttpQueryParameter(string name)
    {
      _name = name;
    }

    internal static void EnumFields(object request, Action<string, string> action)
    {
      foreach (var field in request.GetType().GetProperties(BindingFlags.Instance | BindingFlags.Public))
      {
        var name = field.Name.ToLower();
        var attr = field.GetCustomAttribute<HttpQueryParameter>();
        if (attr != null && !string.IsNullOrEmpty(attr.Name))
        {
          name = attr.Name;
        }
        var val = field.GetValue(request);
        var vs = val.ToString();
        if (!string.IsNullOrEmpty(vs))
        {
          action(name, vs);
        }
      }
    }
    internal static string EncodeFieldQueryString(object request)
    {
      var queryparams = new Dictionary<string, string>();
      EnumFields(request, queryparams.Add);
      var vals = queryparams.Select(p => String.Format("{0}={1}", p.Key, Uri.EscapeDataString(p.Value))).ToList();
      return string.Join("&", vals);
    }

    internal static string EncodeRelativePath(object request, string pattern)
    {
      var pathparams = new Dictionary<string, string>();
      EnumFields(request, pathparams.Add);
      foreach (var p in pathparams)
      {
        var s = "{" + p.Key + "}";
        var v = Uri.EscapeDataString(p.Value);
        pattern = pattern.Replace(s, v);
      }
      return pattern;
    }
  }
}
