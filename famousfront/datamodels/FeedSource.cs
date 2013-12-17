using System.Runtime.Serialization;

namespace famousfront.datamodels
{
  [DataContract]
  internal class FeedSource
  {
    [DataMember(EmitDefaultValue = false)]
    public string name
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public string uri
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public string local
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public uint period
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public ulong deadline
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public uint type
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public int subscribe_state
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public bool enable_proxy
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public ulong update
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public string website
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public string[] tags
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public string[] categories
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public int unreaded
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public FeedMedia media
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public string description
    {
      get;set;
    }
    [DataMember(EmitDefaultValue = false)]
    public string logo
    {
      get;
      set;
    }
  }
}